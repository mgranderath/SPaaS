package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

// ChangePassword allows for changing the password
func (provider *serviceProvider) ChangePassword(c echo.Context) error {
	newPassword := c.FormValue("password")

	err := provider.authService.ChangePassword(newPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

// Authorize is the endpoint for login
func (provider *serviceProvider) Authorize(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	token, err := provider.authService.Authorize(username, password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
