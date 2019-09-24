package routing

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mgranderath/SPaaS/server/di"
	"github.com/mgranderath/SPaaS/server/handlers"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo, p di.Provider) {
	serviceProvider := handlers.NewServiceProvider(p.GetConfigRepository(), p.GetDockerClient())
	e.POST("/login", serviceProvider.Authorize)
	e.File("/", "static/index.html")

	secret := []byte(p.GetConfigRepository().Config.GetString("secret"))
	g := e.Group("")
	g.Use(middleware.JWT(secret))
	g.POST("/change-password", serviceProvider.ChangePassword)

	r := e.Group("/api/app")
	r.Use(middleware.JWT(secret))
	r.GET("", serviceProvider.GetApplications)
	r.GET("/:name", serviceProvider.GetApplication)
	r.POST("/:name", serviceProvider.CreateApplication)
	r.DELETE("/:name", serviceProvider.DeleteApplication)
	r.POST("/:name/start", serviceProvider.StartApplication)
	r.POST("/:name/stop", serviceProvider.StopApplication)
	r.POST("/:name/deploy", serviceProvider.DeployApplication)
	r.GET("/:name/logs", serviceProvider.GetApplicationLogs)
}
