package routing

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mgranderath/SPaaS/server/auth"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/mgranderath/SPaaS/server/service/app"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo, ctx *model.AppDp) {
	authService := auth.NewAuthService(ctx)
	e.POST("/login", authService.Login)
	e.File("/", "static/index.html")

	secret := []byte(ctx.ConfigStore.Config.GetString("secret"))
	g := e.Group("")
	g.Use(middleware.JWT(secret))
	g.POST("/change-password", authService.ChangePassword)

	appService := app.NewAppService(ctx)
	r := e.Group("/api/app")
	r.Use(middleware.JWT(secret))
	r.GET("", appService.GetApplications)
	r.GET("/:name", appService.GetApplication)
	r.POST("/:name", appService.CreateApplication)
	r.DELETE("/:name", appService.DeleteApplication)
	r.POST("/:name/start", appService.StartApplication)
	r.POST("/:name/stop", appService.StopApplication)
	r.POST("/:name/deploy", appService.DeployApplication)
	r.GET("/:name/logs", appService.GetLogs)
}
