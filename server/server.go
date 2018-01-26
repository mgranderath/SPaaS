package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/magrandera/PiaaS/server/docker"

	"github.com/magrandera/PiaaS/server/models"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/routing"
)

func initServer() error {
	dock, err := docker.New()
	if err != nil {
		return err
	}
	list, err := dock.ListContainers()
	if err != nil {
		return err
	}
	proxy := false
	for _, container := range list {
		for _, element := range container.Names {
			if element == "/pi-nginx-proxy" {
				proxy = true
			}
		}
	}
	if !proxy {
		fmt.Println("Nginx Proxy not installed! Installing now!")
		reader, err := dock.Cli.ImagePull(dock.Ctx, "jwilder/nginx-proxy", types.ImagePullOptions{})
		defer reader.Close()
		if err != nil {
			return err
		}
		_, _ = ioutil.ReadAll(reader)
		_, err = dock.Cli.ContainerCreate(dock.Ctx, &container.Config{
			Image: "jwilder/nginx-proxy",
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		}, &container.HostConfig{
			Binds: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "80",
					},
				},
			}}, nil, "pi-nginx-proxy")
		if err != nil {
			return err
		}
		if err = dock.StartContainer("pi-nginx-proxy"); err != nil {
			return err
		}
		fmt.Println("Nginx Proxy was installed!")
	}
	return nil
}

// StartServer : start the PiaaS server
func StartServer() {
	router := mux.NewRouter()
	routing.SetupRouting(router)
	models.InitDB()
	err := initServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(http.ListenAndServe(":5000", router))
}
