package services

import (
	"errors"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
)

func (a *AppService) Start(name string, messages model.StatusChannel) {
	app := model.NewApplication(name)
	if !app.Exists() {
		messages.SendError(errors.New("Does not exist"))
		close(messages)
		return
	}
	messages.SendInfo("Starting application")
	if err := a.DockerClient.StartContainer(common.SpaasName(name)); err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	messages.SendSuccess("Starting application")
	close(messages)
}
