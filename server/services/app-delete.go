package services

import (
	"errors"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
	"os"
)

func (a *AppService) Delete(name string, messages model.StatusChannel) {
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
	_ = a.DockerClient.RemoveContainer(common.SpaasName(name))
	messages.SendSuccess("Removing docker container")
	messages.SendInfo("Removing docker image")
	_, _ = a.DockerClient.RemoveImage(common.SpaasName(name))
	messages.SendSuccess("Removing docker image")
	close(messages)
}
