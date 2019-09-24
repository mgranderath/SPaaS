package di

import (
	"github.com/mgranderath/SPaaS/config"
	"github.com/mgranderath/SPaaS/server/docker"
)

type Provider interface {
	GetDockerClient() *docker.Docker
	GetConfigRepository() *config.Store
}

type provider struct {
	dockerClient     *docker.Docker
	configRepository *config.Store
}

func NewProvider() Provider {
	configRepository := config.New()
	dockerClient := docker.InitDocker()
	return &provider{
		dockerClient,
		configRepository,
	}
}
func (p *provider) GetDockerClient() *docker.Docker {
	return p.dockerClient
}

func (p *provider) GetConfigRepository() *config.Store {
	return p.configRepository
}
