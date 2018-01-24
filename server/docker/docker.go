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
	ctx context.Context
	cli *client.Client
}

// New : creates a new connection to the docker socket
func New() (Docker, error) {
	dock := Docker{}
	dock.ctx = context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return Docker{}, err
	}
	dock.cli = cli
	return dock, nil
}

// BuildImage : builds a docker image from a tar file
func (dock Docker) BuildImage(tarfile *os.File, name string) (types.ImageBuildResponse, error) {
	imageBuildResponse, err := dock.cli.ImageBuild(
		dock.ctx,
		tarfile,
		types.ImageBuildOptions{
			Tags:       []string{"pi-" + name},
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
	_, err := dock.cli.ImageRemove(dock.ctx, "pi-"+name, types.ImageRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such image") {
		return err
	}
	return nil
}

// BuildContainer : builds a docker container
func (dock Docker) BuildContainer(name string, port string) (container.ContainerCreateCreatedBody, error) {
	response, err := dock.cli.ContainerCreate(dock.ctx, &container.Config{
		Image: "pi-" + name,
		Env:   []string{"VIRTUAL_HOST=" + name + ".granderath.tech"},
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
		}}, nil, "pi-"+name)
	if err != nil {
		return response, err
	}
	return response, nil
}

// RemoveContainer : removes a docker container
func (dock Docker) RemoveContainer(name string) error {
	err := dock.cli.ContainerRemove(dock.ctx, "pi-"+name, types.ContainerRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		return err
	}
	return nil
}

// StartContainer : starts a docker container
func (dock Docker) StartContainer(name string) error {
	err := dock.cli.ContainerStart(dock.ctx, "pi-"+name, types.ContainerStartOptions{})
	return err
}

// StopContainer : stops a docker container
func (dock Docker) StopContainer(name string) error {
	err := dock.cli.ContainerStop(dock.ctx, "pi-"+name, nil)
	if err != nil {
		return err
	}
	return nil
}
