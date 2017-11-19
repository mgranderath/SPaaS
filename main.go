package main

import (
	"os"

	"github.com/magrandera/PiaaS/command"
	"github.com/magrandera/PiaaS/models"
	"github.com/magrandera/PiaaS/server"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "PiaaS"
	app.Usage = "A Heroku for the Raspberry Pi"
	models.InitDB()
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "start the PiaaS server",
			Action: func(c *cli.Context) error {
				server.StartServer()
				return nil
			},
		},
		{
			Name:  "app",
			Usage: "options for applications",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all applications",
					Action: func(c *cli.Context) error {
						command.ListApplications()
						return nil
					},
				},
				{
					Name:  "add",
					Usage: "add an applications",
					Action: func(c *cli.Context) error {
						command.CreateApplication(c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
