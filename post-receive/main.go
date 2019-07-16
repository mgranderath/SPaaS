package main

import (
	client "docker.io/go-docker"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	applicationName := os.Args[1]
	docker, err := client.NewEnvClient()
	if err != nil {
		log.Panic("Could not connect to docker container")
	}
}
