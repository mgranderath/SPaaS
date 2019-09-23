package main

import (
	"flag"
	"fmt"
	"github.com/mgranderath/SPaaS/server/model"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/server/routing"
	"github.com/mgranderath/SPaaS/server/service/app"
)

var appDp *model.AppDp

func initialize(e *echo.Echo) {
	appDp = model.NewAppDp()
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e, appDp)
	routing.InitReverseProxy(appDp)
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
		appService := app.NewAppService(appDp)
		go appService.Deploy(appName, messages)
		for elem := range messages {
			fmt.Println(strings.ToUpper(elem.Type) + ": " + elem.Message)
		}
		return
	}
	e.Logger.Fatal(e.Start(":8080"))
}
