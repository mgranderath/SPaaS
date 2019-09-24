package services

import (
	"errors"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
)

func (a *AppService) Stop(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	if !app.Exists() {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	messages.SendInfo("Stopping application")
	if err := a.DockerClient.StopContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Stopping application")
	close(messages)
}
