package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
)

func stop(name string, messages chan<- Application) {
	messages <- Application{
		Type:    "info",
		Message: "Stopping application",
	}
	if err := StopContainer(common.SpaasName(name)); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Stopping application",
	}
	close(messages)
}

// StopApplication starts an application
func StopApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	messages := make(chan Application)
	go stop(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			return c.JSON(http.StatusInternalServerError, Application{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
