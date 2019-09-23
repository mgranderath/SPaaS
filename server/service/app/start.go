package app

import (
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/server/model"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mgranderath/SPaaS/common"
)

func (appService *AppService) start(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	if !app.Exists() {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	messages.SendInfo("Starting application")
	if err := appService.Docker.StartContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Starting application")
	close(messages)
}

// StartApplication starts an application
func (app *AppService) StartApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being started", name)
	messages := make(chan model.Status)
	go app.start(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' start failed with: %v", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
