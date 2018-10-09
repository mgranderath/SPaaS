package controller

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	git "gopkg.in/src-d/go-git.v4"
)

func createDeployDir(deployPath string, messages chan<- Application) error {
	// Creating directory
	messages <- Application{
		Type:    "info",
		Message: "Creating directories",
	}
	if err := os.RemoveAll(deployPath); err != nil {
		return err
	}
	err := os.MkdirAll(deployPath, os.ModePerm)
	if err != nil {
		return err
	}
	messages <- Application{
		Type:    "success",
		Message: "Creating directories",
	}
	return nil
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
	err := createDeployDir(deployPath, messages)
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
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
	response, err := BuildImage(f, common.SpaasName(name))
	if err != nil {
		messages <- Application{
			Type:    "error",
			Message: err.Error(),
		}
		close(messages)
		return
	}
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
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
		"traefik.backend": common.SpaasName(name),
		"traefik.enable":  "true",
		"traefik.port":    "80",
	}
	if config.Cfg.Config.GetBool("useDomain") {
		labels["traefik.frontend.rule"] = "Host:" + name + "." + config.Cfg.Config.GetString("domain")
	} else {
		labels["traefik.frontend.rule"] = "PathPrefixStrip:/spaas/" + name
	}
	_, err = CreateContainer(
		container.Config{
			Image: common.SpaasName(name) + ":latest",
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
			Labels: labels,
		}, container.HostConfig{}, network.NetworkingConfig{}, common.SpaasName(name))
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
