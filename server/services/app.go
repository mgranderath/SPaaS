package services

import (
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/docker"
	"path/filepath"

	"github.com/mgranderath/SPaaS/common"
)

type AppService struct {
	ConfigRespository *config.Store
	DockerClient      *docker.Docker
}

func NewAppService(configRepository *config.Store, dockerClient *docker.Docker) *AppService {
	return &AppService{
		ConfigRespository: configRepository,
		DockerClient:      dockerClient,
	}
}

var basePath = filepath.Join(common.HomeDir(), ".spaas")
