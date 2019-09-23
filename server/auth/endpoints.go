package auth

import (
	"github.com/mgranderath/SPaaS/server/model"
	"net/http"
	"time"

	"github.com/mgranderath/SPaaS/common"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/config"
)

type AuthService struct {
	Config *config.Store
}

func NewAuthService(ctx *model.AppDp) *AuthService {
	return &AuthService{
		Config: ctx.ConfigStore,
	}
}

// ChangePassword allows for changing the password
func (service *AuthService) ChangePassword(c echo.Context) error {
	newPassword := c.FormValue("password")
	if len(newPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password has to be at leat 8 characters",
		})
	}
	hashedPassword := common.HashPassword(newPassword)
	service.Config.Config.Set("password", hashedPassword)
	service.Config.Save()
	return c.NoContent(http.StatusOK)
}

// Login is the endpoint for login
func (service *AuthService) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == service.Config.Config.GetString("username") && common.CheckPasswordHash(password,
		service.Config.Config.GetString("password")) {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["admin"] = true
		claims["created"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(service.Config.Config.GetString("secret")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

// GetToken generates a token for internal request use
func (service *AuthService) GetToken() (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "spaas"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(service.Config.Config.GetString("secret")))
	if err != nil {
		return "", err
	}
	return t, nil
}
