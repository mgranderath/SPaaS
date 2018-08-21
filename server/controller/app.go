package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/magrandera/SPaaS/common"
)

// Application stores information about the application
type Application struct {
	Type     string     `json:"type"`
	Message  string     `json:"message"`
	Extended []KeyValue `json:"extended"`
}

// KeyValue holds extra information of a message
type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var basePath = filepath.Join(common.HomeDir(), ".spaas")

func create(name string, messages chan<- Application) {
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	// Check if app already exists
	if common.Exists(appPath) {
		messages <- Application{
			Type:    "error",
			Message: "Already exists",
		}
		close(messages)
		return
	}
	// Create Directories
	messages <- Application{
		Type:    "info",
		Message: "Creating directories",
	}
	err := os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating directories",
	}
	// Initialize the git repository
	messages <- Application{
		Type:    "info",
		Message: "Creating git repo",
	}
	if _, err := git.PlainInit(repoPath, true); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating git repo",
	}
	// Create git post-receive hook
	messages <- Application{
		Type:    "info",
		Message: "Creating git receive hook",
	}
	err = os.MkdirAll(filepath.Join(repoPath, "hooks"), os.ModePerm)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	file, err := os.Create(filepath.Join(repoPath, "hooks", "post-receive"))
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\necho \"Starting deploy!\"\n"+
		"curl --request POST 'http://127.0.0.1:5000/api/v1/app/%s/deploy' |"+
		"python -c 'import json,sys;obj=json.load(sys.stdin);print \"Successfully deployed!\";print \"Port: \"+obj[\"port\"]'\n", name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(repoPath, "hooks", "post-receive"), 0755)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating git receive hook",
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating app",
		Extended: []KeyValue{
			{Key: "RepoPath", Value: repoPath},
		},
	}
	close(messages)
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
	// Remove directories
	messages <- Application{
		Type:    "info",
		Message: "Removing directories",
	}
	if err := os.RemoveAll(appPath); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Removing directories",
	}
	close(messages)
}

func deploy(name string, messages chan<- Application) {
	appPath := filepath.Join(basePath, "applications", name)
	deployPath := filepath.Join(appPath, "deploy")
	repoPath := filepath.Join(appPath, "repo")
	if !common.Exists(appPath) {
		messages <- Application{
			Type:    "error",
			Message: "Does not exist",
		}
		close(messages)
		return
	}
	// Creating directory
	messages <- Application{
		Type:    "info",
		Message: "Creating directories",
	}
	err := os.MkdirAll(deployPath, os.ModePerm)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating directories",
	}
	// Clone repository
	messages <- Application{
		Type:    "info",
		Message: "Cloning repo",
	}
	_, err = git.PlainClone(deployPath, false, &git.CloneOptions{
		URL: repoPath,
	})
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
}

// CreateApplication creates a new application
func CreateApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	messages := make(chan Application)
	go create(name, messages)
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

// DeployApplication deploys an application
func DeployApplication(c echo.Context) error {
	name := c.Param("name")
	messages := make(chan Application)
	go deploy(name, messages)
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
