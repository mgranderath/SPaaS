package controller

import (
	"bufio"
	"github.com/mgranderath/SPaaS/server/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"

	"github.com/mgranderath/SPaaS/common"
)

var basePath = filepath.Join(common.HomeDir(), ".spaas")

// GetLogs streams the log of a container
func GetLogs(c echo.Context) error {
	name := c.Param("name")
	resp, err := ContainerLogs(common.SpaasName(name))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
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
		if err := common.EncodeJSONAndFlush(c, model.Status{
			Type:    "info",
			Message: line,
		}); err != nil {
			return c.JSON(http.StatusInternalServerError, model.Status{
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
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, container); err != nil {
		return c.JSON(http.StatusInternalServerError, model.Status{
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
		return c.JSON(http.StatusInternalServerError, model.Status{
			Type:    "error",
			Message: err.Error(),
		})
	}
	for _, f := range files {
		if err := common.EncodeJSONAndFlush(c, model.Status{
			Type:    "info",
			Message: f.Name(),
		}); err != nil {
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
