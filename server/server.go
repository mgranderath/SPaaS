package main

import (
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/server/routing"
)

func main() {
	e := echo.New()
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
