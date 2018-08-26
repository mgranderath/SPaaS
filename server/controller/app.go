package controller

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	client "github.com/fsouza/go-dockerclient"
	"github.com/labstack/echo"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/auth"
	"github.com/magrandera/SPaaS/server/hook"
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
	token, err := auth.GetToken()
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	postReceive, err := hook.CreatePostReceive(name, token)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	_, err = file.WriteString(postReceive)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
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
	_ = RemoveImage(common.SpaasName(name))
	messages <- Application{
		Type:    "success",
		Message: "Removing docker image",
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
	if err := os.RemoveAll(deployPath); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
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
	messages <- Application{
		Type:    "success",
		Message: "Cloning repo",
	}
	messages <- Application{
		Type:    "info",
		Message: "Detecting run command",
	}
	dockerfile := config.Dockerfile{}
	v, err := config.ReadConfig(filepath.Join(deployPath, "spaas.json"), map[string]interface{}{})
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	if !v.InConfig("start") {
		messages <- Application{
			Type:    "error",
			Message: "No start in spaas.json in project",
		}
		close(messages)
		return
	}
	dockerfile.Command = strings.Fields(v.GetString("start"))
	messages <- Application{
		Type:    "success",
		Message: "Detecting run command",
		Extended: []KeyValue{
			{Key: "Cmd", Value: v.GetString("start")},
		},
	}
	messages <- Application{
		Type:    "info",
		Message: "Detecting app type",
	}
	if common.Exists(filepath.Join(deployPath, "requirements.txt")) {
		dockerfile.Type = "python"
	} else if common.Exists(filepath.Join(deployPath, "package.json")) {
		dockerfile.Type = "nodejs"
	} else if common.Exists(filepath.Join(deployPath, "Gemfile")) {
		dockerfile.Type = "ruby"
	} else {
		messages <- Application{
			Type:    "error",
			Message: "Could not detect type of application",
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Detecting app type",
		Extended: []KeyValue{
			{Key: "Type", Value: dockerfile.Type},
		},
	}
	messages <- Application{
		Type:    "info",
		Message: "Packaging app",
	}
	if err := config.CreateDockerfile(dockerfile, appPath); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	cmd := exec.Command("tar", "cvf", "../package.tar", ".")
	cmd.Dir = deployPath + "/"
	_, err = cmd.Output()
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
		Message: "Packaging app",
	}
	messages <- Application{
		Type:    "info",
		Message: "Building image",
	}
	f, err := os.Open(filepath.Join(appPath, "package.tar"))
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	defer f.Close()
	if err := BuildImage(f, common.SpaasName(name)); err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	messages <- Application{
		Type:    "success",
		Message: "Building image",
	}
	messages <- Application{
		Type:    "info",
		Message: "Building container",
	}
	_ = RemoveContainer(common.SpaasName(name))
	labels := map[string]string{
		"traefik.backend":       common.SpaasName(name),
		"traefik.frontend.rule": "Host:" + name + ".granderath.tech",
		"traefik.enable":        "true",
		"traefik.port":          "80",
	}
	_, err = CreateContainer(client.CreateContainerOptions{
		Name: common.SpaasName(name),
		Config: &client.Config{
			Image: common.SpaasName(name) + ":latest",
			ExposedPorts: map[client.Port]struct{}{
				"80/tcp": struct{}{},
			},
			Labels: labels,
		},
		HostConfig: &client.HostConfig{},
	})
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
		Message: "Building container",
	}
	messages <- Application{
		Type:    "info",
		Message: "Starting container",
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
		Message: "Starting container",
	}
	close(messages)
}

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
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
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
