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
	e.File("/", "static/index.html")

	g := e.Group("")
	g.Use(middleware.JWT([]byte(config.Cfg.Config.GetString("secret"))))
	g.POST("/change-password", auth.ChangePassword)

	r := e.Group("/api/app")
	r.Use(middleware.JWT([]byte(config.Cfg.Config.GetString("secret"))))
	r.GET("", controller.GetApplications)
	r.GET("/:name", controller.GetApplication)
	r.POST("/:name", controller.CreateApplication)
	r.DELETE("/:name", controller.DeleteApplication)
	r.POST("/:name/start", controller.StartApplication)
	r.POST("/:name/stop", controller.StopApplication)
	r.POST("/:name/deploy", controller.DeployApplication)
	r.GET("/:name/logs", controller.GetLogs)
}
