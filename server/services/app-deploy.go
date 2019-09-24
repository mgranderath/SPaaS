package services

import (
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func packageApplication(app *model.Application) error {
	if err := app.Build(); err != nil {
		return err
	}
	cmd := exec.Command("tar", "cvf", "../package.tar", ".")
	cmd.Dir = app.DeployPath + "/"
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func (a *AppService) buildImage(app *model.Application) error {
	f, err := os.Open(filepath.Join(app.Path, "package.tar"))
	if err != nil {
		return err
	}
	defer f.Close()
	response, err := a.DockerClient.BuildImage(f, common.SpaasName(app.Name))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	return nil
}

func (a *AppService) buildContainer(app *model.Application) error {
	name := app.Name
	_ = a.DockerClient.RemoveContainer(common.SpaasName(name))
	labels := map[string]string{
		"traefik.backend": common.SpaasName(name),
		"traefik.enable":  "true",
		"traefik.port":    "80",
	}
	if a.ConfigRespository.Config.GetBool("useDomain") {
		labels["traefik.frontend.rule"] =
			fmt.Sprintf("Host:%s.%s", name, a.ConfigRespository.Config.GetString("domain"))
	} else {
		labels["traefik.frontend.rule"] =
			fmt.Sprintf("PathPrefixStrip:/spaas/%s", name)
	}
	_, err := a.DockerClient.CreateContainer(
		container.Config{
			Image: common.SpaasName(name) + ":latest",
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
			Env:    []string{"PORT=80"},
			Labels: labels,
			Tty:    true,
		}, container.HostConfig{}, network.NetworkingConfig{}, common.SpaasName(name))
	return err
}

func (a *AppService) Deploy(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	if !app.Exists() {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	messages.SendInfo("Creating directories")
	if err := app.ResetDeployDir(); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Creating directories")
	// Clone repository
	messages.SendInfo("Cloning repository")
	_, err := git.PlainClone(app.DeployPath, false, &git.CloneOptions{
		URL: app.RepositoryPath,
	})
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Cloning repository")
	messages.SendInfo("Detecting run command")
	if err := app.DetectStartCommand(); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages <- model.Status{
		Type:    "success",
		Message: "Detecting run command",
	}
	messages.SendInfo("Detecting app type")
	if app.DetectType() == model.Undefined {
		messages.SendError(errors.New("Could not detect type of application"))
		close(messages)
		return
	}
	messages <- model.Status{
		Type:    "success",
		Message: "Detecting app type",
	}
	messages.SendInfo("Packaging application")
	err = packageApplication(app)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Packaging application")
	messages.SendInfo("Building image")
	err = a.buildImage(app)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Building image")
	messages.SendInfo("Building container")
	err = a.buildContainer(app)
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Building container")
	messages.SendInfo("Starting container")
	if err := a.DockerClient.StartContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Starting container")
	close(messages)
}
