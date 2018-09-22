package routing

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// GlobalMiddleware applies the global middleware to the router
func GlobalMiddleware(e *echo.Echo) {
	e.Static("/static", "static")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
}
