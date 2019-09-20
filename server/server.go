package main

import (
	"flag"
	"fmt"
	"github.com/mgranderath/SPaaS/server/model"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/controller"
	"github.com/mgranderath/SPaaS/server/routing"
)

func initialize(e *echo.Echo) {
	config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	if err := config.Save(); err != nil {
		fmt.Println(err.Error())
	}
	config.Cfg.Config.WatchConfig()
	config.Cfg.Config.OnConfigChange(func(_ fsnotify.Event) {
		fmt.Println("Config file changed")
	})
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e)
	controller.InitDocker()
	routing.InitReverseProxy()
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
		go controller.Deploy(appName, messages)
		for elem := range messages {
			fmt.Println(strings.ToUpper(elem.Type) + ": " + elem.Message)
		}
		return
	}
	e.Logger.Fatal(e.Start(":8080"))
}
