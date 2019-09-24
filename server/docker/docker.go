package docker

import (
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
)

type DockerClient interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, buildContext io.Reader,
		options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImageRemove(ctx context.Context, imageID string,
		options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
}

// DockerClient holds the connection information for the docker instance
type Docker struct {
	Ctx context.Context
	Cli DockerClient
}

var newDockerClient = func() (DockerClient, error) {
	return client.NewEnvClient()
}

// InitDocker initializes the docker instance
func InitDocker() *Docker {
	dock := Docker{}
	dock.Ctx = context.Background()
	Cli, err := newDockerClient()
	if err != nil {
		log.Panic("Could not connect to DockerClient")
	}
	dock.Cli = Cli
	return &dock
}

// ListContainers retrieves a list of all containers running on the system
func (dock *Docker) ListContainers() ([]types.Container, error) {
	list, err := dock.Cli.ContainerList(dock.Ctx, types.ContainerListOptions{
		Quiet: true,
		All:   true,
	})
	return list, err
}

// PullImage pulls an image from the docker registry
func (dock *Docker) PullImage(name string) error {
	reader, err := dock.Cli.ImagePull(dock.Ctx, name, types.ImagePullOptions{})
	defer reader.Close()
	_, err = ioutil.ReadAll(reader)
	return err
}

// CreateContainer creates a container
func (dock *Docker) CreateContainer(
	containerConfig container.Config,
	hostConfig container.HostConfig,
	networkConfig network.NetworkingConfig,
	containerName string,
) (container.ContainerCreateCreatedBody, error) {
	return dock.Cli.ContainerCreate(dock.Ctx, &containerConfig, &hostConfig, &networkConfig, containerName)
}

// StartContainer starts the container with id
func (dock *Docker) StartContainer(id string) error {
	return dock.Cli.ContainerStart(dock.Ctx, id, types.ContainerStartOptions{})
}

// StopContainer stops the container with id
func (dock *Docker) StopContainer(id string) error {
	zero := (0 * time.Microsecond)
	return dock.Cli.ContainerStop(dock.Ctx, id, &zero)
}

// BuildImage builds an image from a tar stream
func (dock *Docker) BuildImage(tarfile *os.File, name string) (types.ImageBuildResponse, error) {
	return dock.Cli.ImageBuild(dock.Ctx, tarfile, types.ImageBuildOptions{
		Tags: []string{name},
	})
}

// RemoveContainer removes an container
func (dock *Docker) RemoveContainer(name string) error {
	return dock.Cli.ContainerRemove(dock.Ctx, name, types.ContainerRemoveOptions{
		Force: true,
	})
}

// RemoveImage removes an image
func (dock *Docker) RemoveImage(name string) ([]types.ImageDeleteResponseItem, error) {
	return dock.Cli.ImageRemove(dock.Ctx, name, types.ImageRemoveOptions{
		Force: true,
	})
}

// InspectContainer inspects a container
func (dock *Docker) InspectContainer(name string) (types.ContainerJSON, error) {
	return dock.Cli.ContainerInspect(dock.Ctx, name)
}

// ContainerLogs returns log of a container
func (dock *Docker) ContainerLogs(name string) (io.ReadCloser, error) {
	return dock.Cli.ContainerLogs(dock.Ctx, name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "100",
	})
}
