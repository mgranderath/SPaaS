package auth

import (
	"net/http"
	"time"

	"github.com/magrandera/SPaaS/common"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/config"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == config.Cfg.Config.GetString("username") && common.CheckPasswordHash(password, config.Cfg.Config.GetString("password")) {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.Cfg.Config.GetString("secret")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}