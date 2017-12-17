package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/magrandera/PiaaS/server/models"
	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

type configFile struct {
	Server string
}

func getConf() configFile {
	c := configFile{}
	config := filepath.Join(models.GetHomeFolder(), ".config", "piaas", "config")
	file, err := ioutil.ReadFile(config)
	if err != nil {
		return c
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return c
	}
	return c
}

func main() {
	app := cli.NewApp()
	app.Name = "PiaaS"
	app.Usage = "A Heroku for the Raspberry Pi"
	models.InitDB()
	app.Commands = []cli.Command{
		{
			Name:  "list",
			Usage: "list all applications",
			Action: func(c *cli.Context) error {
				config := getConf()
				if config.Server == "" {
					models.PrintErr(os.Stdout, "Config file has not been created. Run setup")
					return nil
				}
				res, err := http.Get("http://" + config.Server + "/api/v1/app")
				if err != nil {
					models.PrintErr(os.Stdout, err)
					return nil
				}
				if res.StatusCode != http.StatusOK {
					models.PrintErr(os.Stdout, "response")
					return nil
				}
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				apps := []models.Application{}
				json.Unmarshal(bodyBytes, &apps)
				for _, app := range apps {
					fmt.Println(app.Name)
				}
				return nil
			},
		},
		{
			Name:  "setup",
			Usage: "create the configuration file",
			Action: func(c *cli.Context) error {
				newConfig := configFile{}
				scanner := bufio.NewScanner(os.Stdin)
				models.PrintInfo(os.Stdout, "Enter the server ip/url")
				fmt.Print("server: ")
				scanner.Scan()
				newConfig.Server = scanner.Text()
				d, err := yaml.Marshal(&newConfig)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				err = os.MkdirAll(filepath.Join(models.GetHomeFolder(), ".config", "piaas"), os.ModePerm)
				if err != nil {
					return nil
				}
				config := filepath.Join(models.GetHomeFolder(), ".config", "piaas", "config")
				_, err = os.Stat(config)
				if os.IsNotExist(err) {
					var file, err = os.Create(config)
					if err != nil {
						return nil
					}
					defer file.Close()
				}
				file, err := os.OpenFile(config, os.O_RDWR, 0644)
				if err != nil {
					return nil
				}
				file.Write(d)
				defer file.Close()
				return nil
			},
		},
		{
			Name:  "add",
			Usage: "add an application",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove an application",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "deploy",
			Usage: "deploy an application",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "start an application",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "stop an application",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "log of application",
			Action: func(c *cli.Context) error {
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "tail", Usage: "Tail the log"},
			},
		},
	}

	app.Run(os.Args)
}
