package controller

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
	"net/http"
	"os"
)

func delete(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	if !common.Exists(app.Path) {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	// Remove directories
	messages.SendInfo("Removing directories")
	if err := os.RemoveAll(app.Path); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Removing directories")
	messages.SendInfo("Removing docker container")
	_ = RemoveContainer(common.SpaasName(name))
	messages.SendSuccess("Removing docker container")
	messages.SendInfo("Removing docker image")
	_, _ = RemoveImage(common.SpaasName(name))
	messages.SendSuccess("Removing docker image")
	close(messages)
}

// DeleteApplication deletes the application
func DeleteApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being deleted\n", name)
	messages := make(chan model.Status)
	go delete(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' deletion failed with: %v\n", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
