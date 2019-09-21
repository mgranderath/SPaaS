package main

import (
	"context"
	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	applicationName := os.Args[1]
	docker, err := client.NewEnvClient()
	if err != nil || docker == nil {
		log.Panic("Could not connect to docker container")
		os.Exit(1)
	}
	response, err := docker.ContainerExecCreate(context.Background(), "spaas", types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/app/SPaaS_server", "-deploy", applicationName},
	})
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	execId := response.ID
	hijackedResponse, err := docker.ContainerExecAttach(context.Background(), execId, types.ExecConfig{})
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	defer hijackedResponse.Close()
	// Using stdcopy because container without a tty prepend output with an extra byte
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, hijackedResponse.Reader)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}
