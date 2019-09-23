package model

import (
	"fmt"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/docker"
	"path/filepath"
)

type AppDp struct {
	ConfigStore *config.Store
	Docker      *docker.Docker
}

func NewAppDp() *AppDp {
	Config := config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	if err := Config.Save(); err != nil {
		fmt.Println(err.Error())
	}
	Docker := docker.InitDocker()
	return &AppDp{
		ConfigStore: Config,
		Docker:      Docker,
	}
}
