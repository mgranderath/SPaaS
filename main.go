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
					Usage: "add an application",
					Action: func(c *cli.Context) error {
						command.CreateApplication(c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an application",
					Action: func(c *cli.Context) error {
						command.DeleteApplication(c.Args().First())
						return nil
					},
				},
				{
					Name:  "deploy",
					Usage: "deploy an application",
					Action: func(c *cli.Context) error {
						command.DeployApplication(c.Args().First())
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "start an application",
					Action: func(c *cli.Context) error {
						command.StartApplication(c.Args().First())
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "stop an application",
					Action: func(c *cli.Context) error {
						command.StopApplication(c.Args().First())
						return nil
					},
				},
				{
					Name:  "logs",
					Usage: "log of application",
					Action: func(c *cli.Context) error {
						command.LogApplication(c.Args().First(), c.BoolT("tail"))
						return nil
					},
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "tail", Usage: "Tail the log"},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
