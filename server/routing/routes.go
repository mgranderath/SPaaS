package routing

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/auth"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(e *echo.Echo) {
	e.POST("/login", auth.Login)

	r := e.Group("/api")
	r.Use(middleware.JWT([]byte(config.Cfg.Config.GetString("secret"))))
	r.GET("", restricted)
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
