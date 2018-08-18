package routing

import (
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/server/auth"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo) {
	e.POST("/login", auth.Login)
}
