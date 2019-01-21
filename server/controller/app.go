package controller

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"

	"github.com/mgranderath/SPaaS/common"
)

// Application stores information about the application
type Application struct {
	Type     string     `json:"type"`
	Message  string     `json:"message"`
	Extended []KeyValue `json:"extended,omitempty"`
}

// KeyValue holds extra information of a message
type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var basePath = filepath.Join(common.HomeDir(), ".spaas")

// GetLogs streams the log of a container
func GetLogs(c echo.Context) error {
	name := c.Param("name")
	resp, err := ContainerLogs(common.SpaasName(name))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	defer resp.Close()
	rd := bufio.NewReader(resp)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return err
		}
		if err := common.EncodeJSONAndFlush(c, Application{
			Type:    "info",
			Message: line,
		}); err != nil {
			return c.JSON(http.StatusInternalServerError, Application{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}

// GetApplication returns a current application
func GetApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	container, err := InspectContainer(common.SpaasName(name))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, container); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	return nil
}

// GetApplications returns a list of all applications
func GetApplications(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	appPath := filepath.Join(basePath, "applications")
	files, err := ioutil.ReadDir(appPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	for _, f := range files {
		if err := common.EncodeJSONAndFlush(c, Application{
			Type:    "info",
			Message: f.Name(),
		}); err != nil {
			return c.JSON(http.StatusInternalServerError, Application{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
