package controller

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
)

// Docker holds the connection information for the docker instance
type Docker struct {
	Ctx context.Context
	Cli *client.Client
}

var dock Docker

// InitDocker initializes the docker instance
func InitDocker() {
	dock = Docker{}
	dock.Ctx = context.Background()
	Cli, err := client.NewEnvClient()
	if err != nil {
		log.Panic("Could not connect to Docker")
	}
	dock.Cli = Cli
}

// ListContainers retrieves a list of all containers running on the system
func ListContainers() ([]types.Container, error) {
	list, err := dock.Cli.ContainerList(dock.Ctx, types.ContainerListOptions{
		Quiet: true,
		All:   true,
	})
	return list, err
}

// PullImage pulls an image from the docker registry
func PullImage(name string) error {
	reader, err := dock.Cli.ImagePull(dock.Ctx, name, types.ImagePullOptions{})
	defer reader.Close()
	_, err = ioutil.ReadAll(reader)
	return err
}

// CreateContainer creates a container
func CreateContainer(
	containerConfig container.Config,
	hostConfig container.HostConfig,
	networkConfig network.NetworkingConfig,
	containerName string,
) (container.ContainerCreateCreatedBody, error) {
	return dock.Cli.ContainerCreate(dock.Ctx, &containerConfig, &hostConfig, &networkConfig, containerName)
}

// StartContainer starts the container with id
func StartContainer(id string) error {
	return dock.Cli.ContainerStart(dock.Ctx, id, types.ContainerStartOptions{})
}

// StopContainer stops the container with id
func StopContainer(id string) error {
	zero := (0 * time.Microsecond)
	return dock.Cli.ContainerStop(dock.Ctx, id, &zero)
}

// BuildImage builds an image from a tar stream
func BuildImage(tarfile *os.File, name string) (types.ImageBuildResponse, error) {
	return dock.Cli.ImageBuild(dock.Ctx, tarfile, types.ImageBuildOptions{
		Tags: []string{name},
	})
}

// RemoveContainer removes an container
func RemoveContainer(name string) error {
	return dock.Cli.ContainerRemove(dock.Ctx, name, types.ContainerRemoveOptions{
		Force: true,
	})
}

// RemoveImage removes an image
func RemoveImage(name string) ([]types.ImageDeleteResponseItem, error) {
	return dock.Cli.ImageRemove(dock.Ctx, name, types.ImageRemoveOptions{
		Force: true,
	})
}

// InspectContainer inspects a container
func InspectContainer(name string) (types.ContainerJSON, error) {
	return dock.Cli.ContainerInspect(dock.Ctx, name)
}

// ContainerLogs returns log of a container
func ContainerLogs(name string) (io.ReadCloser, error) {
	now := time.Now()
	then := now.Add(time.Minute * -30)
	return dock.Cli.ContainerLogs(dock.Ctx, name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Since:      strconv.FormatInt(then.Unix(), 10),
	})
}
