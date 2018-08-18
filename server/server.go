package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/routing"
)

func main() {
	e := echo.New()
	config.New(common.HomeDir(), ".spaas.json")
	if err := config.Save(); err != nil {
		fmt.Println(err.Error())
	}
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
