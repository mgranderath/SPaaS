package docker

import (
	"context"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Docker : struct that contains the connection information to the docker socket
type Docker struct {
	Ctx context.Context
	Cli *client.Client
}

// New : creates a new connection to the docker socket
func New() (Docker, error) {
	dock := Docker{}
	dock.Ctx = context.Background()
	Cli, err := client.NewEnvClient()
	if err != nil {
		return Docker{}, err
	}
	dock.Cli = Cli
	return dock, nil
}

// BuildImage : builds a docker image from a tar file
func (dock Docker) BuildImage(tarfile *os.File, name string) (types.ImageBuildResponse, error) {
	imageBuildResponse, err := dock.Cli.ImageBuild(
		dock.Ctx,
		tarfile,
		types.ImageBuildOptions{
			Tags:       []string{name},
			Dockerfile: "Dockerfile",
			Remove:     true,
			NoCache:    true})
	if err != nil {
		return imageBuildResponse, err
	}
	return imageBuildResponse, nil
}

// RemoveImage : removes a docker image
func (dock Docker) RemoveImage(name string) error {
	_, err := dock.Cli.ImageRemove(dock.Ctx, name, types.ImageRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such image") {
		return err
	}
	return nil
}

// BuildContainer : builds a docker container
func (dock Docker) BuildContainer(name string, image string, hostname string, port string) (container.ContainerCreateCreatedBody, error) {
	response, err := dock.Cli.ContainerCreate(dock.Ctx, &container.Config{
		Image: name,
		Env:   []string{"VIRTUAL_HOST=" + hostname + ".granderath.tech"},
		ExposedPorts: nat.PortSet{
			"5000/tcp": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"5000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		}}, nil, name)
	if err != nil {
		return response, err
	}
	return response, nil
}

// RemoveContainer : removes a docker container
func (dock Docker) RemoveContainer(name string) error {
	err := dock.Cli.ContainerRemove(dock.Ctx, name, types.ContainerRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		return err
	}
	return nil
}

// StartContainer : starts a docker container
func (dock Docker) StartContainer(name string) error {
	err := dock.Cli.ContainerStart(dock.Ctx, name, types.ContainerStartOptions{})
	return err
}

// StopContainer : stops a docker container
func (dock Docker) StopContainer(name string) error {
	err := dock.Cli.ContainerStop(dock.Ctx, name, nil)
	if err != nil {
		return err
	}
	return nil
}

// ListContainers : returns a list of docker containers
func (dock Docker) ListContainers() ([]types.Container, error) {
	containers, err := dock.Cli.ContainerList(dock.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return containers, err
	}
	return containers, nil
}
