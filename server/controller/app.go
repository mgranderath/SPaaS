package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/magrandera/SPaaS/common"
)

// Application stores information about the application
type Application struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var basePath = filepath.Join(common.HomeDir(), ".spaas")

// CreateApplication creates a new application
func CreateApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	// Check if app already exists
	if common.Exists(appPath) {
		return c.JSON(http.StatusConflict, Application{
			Type:    "error",
			Message: "Already exists",
		})
	}
	// Create Directories
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "info",
		Message: "Creating directories",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	err := os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "success",
		Message: "Creating directories",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	// Initialize the git repository
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "info",
		Message: "Creating git repo",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if _, err := git.PlainInit(repoPath, true); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "success",
		Message: "Creating git repo",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	// Create git pos-receive hook
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "info",
		Message: "Creating git receive hook",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	err = os.MkdirAll(filepath.Join(repoPath, "hooks"), os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	file, err := os.Create(filepath.Join(repoPath, "hooks", "post-receive"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\necho \"Starting deploy!\"\n"+
		"curl --request POST 'http://127.0.0.1:5000/api/v1/app/%s/deploy' |"+
		"python -c 'import json,sys;obj=json.load(sys.stdin);print \"Successfully deployed!\";print \"Port: \"+obj[\"port\"]'\n", name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(repoPath, "hooks", "post-receive"), 0755)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	if err := common.EncodeJSONAndFlush(c, Application{
		Type:    "success",
		Message: "Creating git receive hook",
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Application{
			Type:    "error",
			Message: err.Error(),
		})
	}
	return nil
}

func getApplications(c echo.Context) error {
	return nil
}
