package routing

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/auth"
	"github.com/mgranderath/SPaaS/server/controller"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo) {
	e.POST("/login", auth.Login)
	e.File("/", "static/index.html")

	secret := []byte(config.Cfg.Config.GetString("secret"))
	g := e.Group("")
	g.Use(middleware.JWT(secret))
	g.POST("/change-password", auth.ChangePassword)

	r := e.Group("/api/app")
	r.Use(middleware.JWT(secret))
	r.GET("", controller.GetApplications)
	r.GET("/:name", controller.GetApplication)
	r.POST("/:name", controller.CreateApplication)
	r.DELETE("/:name", controller.DeleteApplication)
	r.POST("/:name/start", controller.StartApplication)
	r.POST("/:name/stop", controller.StopApplication)
	r.POST("/:name/deploy", controller.DeployApplication)
	r.GET("/:name/logs", controller.GetLogs)
}
