package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/common"
)

func removeDirs(appPath string, messages chan<- Application) error {
	// Remove directories
	messages <- Application{
		Type:    "info",
		Message: "Removing directories",
	}
	if err := os.RemoveAll(appPath); err != nil {
		return err
	}
	messages <- Application{
		Type:    "success",
		Message: "Removing directories",
	}
	return nil
}

func removeDockerImageAndContainer(name string, messages chan<- Application) error {
	messages <- Application{
		Type:    "info",
		Message: "Removing docker container",
	}
	_ = RemoveContainer(common.SpaasName(name))
	messages <- Application{
		Type:    "success",
		Message: "Removing docker container",
	}
	messages <- Application{
		Type:    "info",
		Message: "Removing docker image",
	}
	_, _ = RemoveImage(common.SpaasName(name))
	messages <- Application{
		Type:    "success",
		Message: "Removing docker image",
	}
	return nil
}

func delete(name string, messages chan<- Application) {
	appPath := filepath.Join(basePath, "applications", name)
	if !common.Exists(appPath) {
		messages <- Application{
			Type:    "error",
			Message: "Does not exist",
		}
		close(messages)
		return
	}
	err := removeDirs(appPath, messages)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	removeDockerImageAndContainer(name, messages)
	close(messages)
}

// DeleteApplication deletes the application
func DeleteApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	messages := make(chan Application)
	go delete(name, messages)
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
