package controller

import (
	"context"
	"log"

	client "github.com/fsouza/go-dockerclient"
)

// Docker holds the connection information for the docker instance
type Docker struct {
	Ctx context.Context
	Cli *client.Client
}

var dock Docker

func init() {
	dock = Docker{}
	dock.Ctx = context.Background()
	Cli, err := client.NewClientFromEnv()
	if err != nil {
		log.Panic("Could not connect to Docker")
	}
	dock.Cli = Cli
}
