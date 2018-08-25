package controller

import (
	"bytes"
	"context"
	"log"
	"os"

	client "github.com/fsouza/go-dockerclient"
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
	Cli, err := client.NewClientFromEnv()
	if err != nil {
		log.Panic("Could not connect to Docker")
	}
	dock.Cli = Cli
}

// ListContainers retrieves a list of all containers running on the system
func ListContainers() ([]client.APIContainers, error) {
	list, err := dock.Cli.ListContainers(client.ListContainersOptions{
		All: true,
	})
	return list, err
}

// PullImage pulls an image from the docker registry
func PullImage(name string, tag string) error {
	err := dock.Cli.PullImage(client.PullImageOptions{
		Repository: name,
		Tag:        tag,
	}, client.AuthConfiguration{})
	return err
}

// CreateContainer creates a container
func CreateContainer(opts client.CreateContainerOptions) (*client.Container, error) {
	return dock.Cli.CreateContainer(opts)
}

// StartContainer starts the container with id
func StartContainer(id string) error {
	err := dock.Cli.StartContainer(id, &client.HostConfig{})
	return err
}

// StopContainer stops the container with id
func StopContainer(id string) error {
	return dock.Cli.StopContainer(id, 0)
}

// BuildImage builds an image from a tar stream
func BuildImage(tarfile *os.File, name string) error {
	return dock.Cli.BuildImage(client.BuildImageOptions{
		Name:                name,
		ForceRmTmpContainer: true,
		InputStream:         tarfile,
		OutputStream:        bytes.NewBuffer(nil),
		RmTmpContainer:      true,
	})
}

// RemoveContainer removes an container
func RemoveContainer(name string) error {
	return dock.Cli.RemoveContainer(client.RemoveContainerOptions{
		ID:    name,
		Force: true,
	})
}

// RemoveImage removes an image
func RemoveImage(name string) error {
	return dock.Cli.RemoveImage(name)
}
