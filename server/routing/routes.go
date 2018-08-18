package routing

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/auth"
	"github.com/magrandera/SPaaS/server/controller"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo) {
	e.POST("/login", auth.Login)

	r := e.Group("/api/app")
	r.Use(middleware.JWT([]byte(config.Cfg.Config.GetString("secret"))))
	r.POST("/:name", controller.CreateApplication)
}
