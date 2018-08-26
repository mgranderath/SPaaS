package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/common"
)

func start(name string, messages chan<- Application) {
	messages <- Application{
		Type:    "info",
		Message: "Starting application",
	}
	if err := StartContainer(common.SpaasName(name)); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Starting application",
	}
	close(messages)
}

// StartApplication starts an application
func StartApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	messages := make(chan Application)
	go start(name, messages)
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
