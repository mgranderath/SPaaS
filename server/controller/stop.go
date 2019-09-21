package controller

import (
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/server/model"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
)

func stop(name string, messages model.StatusChannel) {
	messages.SendInfo("Stopping application")
	if err := StopContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Stopping application")
	close(messages)
}

// StopApplication starts an application
func StopApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being stopped", name)
	messages := make(chan model.Status)
	go stop(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' stop failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
