package main

import (
	"flag"
	"fmt"
	"github.com/mgranderath/SPaaS/server/di"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/mgranderath/SPaaS/server/services"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/server/routing"
)

func initialize(e *echo.Echo) {
	provider := di.NewProvider()
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e, provider)
	routing.InitReverseProxy(provider)
}

func main() {
	deploy := flag.Bool("deploy", false, "call this to deploy app")
	flag.Parse()
	e := echo.New()
	initialize(e)
	if *deploy {
		// this is for the post-receive hook
		if len(flag.Args()) != 1 {
			fmt.Println("no args passed")
			os.Exit(1)
		}
		appName := flag.Args()[0]
		messages := make(chan model.Status)
		provider := di.NewProvider()
		appService := services.NewAppService(provider.GetConfigRepository(), provider.GetDockerClient())
		go appService.Deploy(appName, messages)
		for elem := range messages {
			fmt.Println(strings.ToUpper(elem.Type) + ": " + elem.Message)
		}
		return
	}
	e.Logger.Fatal(e.Start(":8080"))
}
