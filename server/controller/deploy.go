package controller

import (
	"github.com/labstack/gommon/log"
	buildpacknodejs "github.com/mgranderath/SPaaS/buildpack/nodejs"
	buildpackpython "github.com/mgranderath/SPaaS/buildpack/python"
	buildpackruby "github.com/mgranderath/SPaaS/buildpack/ruby"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/pkg/errors"
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
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"gopkg.in/src-d/go-git.v4"
)

func deploy(name string, messages model.StatusChannel) {
	var (
		appType model.ApplicationType
	)
	appPath := filepath.Join(basePath, "applications", name)
	deployPath := filepath.Join(appPath, "deploy")
	repoPath := filepath.Join(appPath, "repo")
	if !common.Exists(appPath) {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	messages.SendInfo("Creating directories")
	if err := os.RemoveAll(deployPath); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	err := os.MkdirAll(deployPath, os.ModePerm)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Creating directories")
	// Clone repository
	messages.SendInfo("Cloning repository")
	_, err = git.PlainClone(deployPath, false, &git.CloneOptions{
		URL: repoPath,
	})
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Cloning repository")
	messages.SendInfo("Detecting run command")
	dockerfile := config.Dockerfile{}
	v, err := config.ReadConfig(filepath.Join(deployPath, "spaas.json"), map[string]interface{}{})
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	if !v.InConfig("start") {
		messages.SendError(errors.New("No 'start' in spaas.json in project"))
		close(messages)
		return
	}
	dockerfile.Command = strings.Fields(v.GetString("start"))
	messages <- model.Status{
		Type:    "success",
		Message: "Detecting run command",
		Extended: []model.KeyValue{
			{Key: "Cmd", Value: v.GetString("start")},
		},
	}
	messages.SendInfo("Detecting app type")
	if common.Exists(filepath.Join(deployPath, "requirements.txt")) {
		appType = model.Python
	} else if common.Exists(filepath.Join(deployPath, "package.json")) {
		appType = model.Node
	} else if common.Exists(filepath.Join(deployPath, "Gemfile")) {
		appType = model.Ruby
	} else {
		messages.SendError(errors.New("Could not detect type of application"))
		close(messages)
		return
	}
	messages <- model.Status{
		Type:    "success",
		Message: "Detecting app type",
		Extended: []model.KeyValue{
			{Key: "Type", Value: appType.ToString()},
		},
	}
	messages.SendInfo("Packaging application")
	dockerfileConfig := config.Dockerfile{
		Command: dockerfile.Command,
	}

	switch appType {
	case model.Python:
		if err := buildpackpython.Build(appPath, dockerfileConfig); err != nil {
			messages.SendError(err)
			close(messages)
			return
		}
	case model.Node:
		if err := buildpacknodejs.Build(appPath, dockerfileConfig); err != nil {
			messages.SendError(err)
			close(messages)
			return
		}
	case model.Ruby:
		if err := buildpackruby.Build(appPath, dockerfileConfig); err != nil {
			messages.SendError(err)
			close(messages)
			return
		}
	}
	cmd := exec.Command("tar", "cvf", "../package.tar", ".")
	cmd.Dir = deployPath + "/"
	_, err = cmd.Output()
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Packaging application")
	messages.SendInfo("Building image")
	f, err := os.Open(filepath.Join(appPath, "package.tar"))
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	defer f.Close()
	response, err := BuildImage(f, common.SpaasName(name))
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	messages.SendSuccess("Building image")
	messages.SendInfo("Building container")
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
			Env:    []string{"PORT=80"},
			Labels: labels,
		}, container.HostConfig{}, network.NetworkingConfig{}, common.SpaasName(name))
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Building container")
	messages.SendInfo("Starting container")
	if err := StartContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Starting container")
	close(messages)
}

// DeployApplication deploys an application
func DeployApplication(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	name := c.Param("name")
	log.Infof("application '%s' is being deployed\n", name)
	messages := make(chan model.Status)
	go deploy(name, messages)
	for elem := range messages {
		if err := common.EncodeJSONAndFlush(c, elem); err != nil {
			log.Errorf("application '%s' deployment failed with: %v\n", name, err)
			return c.JSON(http.StatusInternalServerError, model.Status{
				Type:    "error",
				Message: err.Error(),
			})
		}
	}
	return nil
}
