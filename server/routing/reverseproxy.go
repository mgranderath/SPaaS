package routing

import (
	"log"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/controller"
)

// InitReverseProxy initializes the reverse proxy
func InitReverseProxy() {
	list, err := controller.ListContainers()
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range list {
		for _, element := range container.Names {
			if element == "/"+common.SpaasName("traefik") {
				return
			}
		}
	}
	log.Println("Traefik is not installed but installing now")
	if err := controller.PullImage("traefik:1.7-alpine"); err != nil {
		log.Fatal(err.Error())
	}
	cmd := []string{
		"--docker", "--docker.watch",
		"--defaultEntryPoints=http",
		"--entryPoints=Name:http Address::80 Compress:off",
		"--docker.domain=granderath.tech",
		"--debug",
		"--logLevel=DEBUG",
	}
	letsencrypt := []string{
		"--acme",
		"--acme.email=" + config.Cfg.Config.GetString("letsencryptEmail"),
		"--acme.storage=/var/acme/acme.json",
		"--acme.httpchallenge.entrypoint=http",
		"--acme.entrypoint=https",
		"--acme.onhostrule=true",
		"--accesslogsfile=/var/acme/access.log",
		"--entryPoints=Name:https Address::443 TLS Compress:off",
		"--entryPoints=Name:http Address::80 Redirect.EntryPoint:https Compress:off",
		"--defaultEntryPoints=https,http",
	}
	if config.Cfg.Config.GetBool("letsencrypt") {
		cmd = append(cmd, letsencrypt...)
	}
<<<<<<< HEAD
	container, err := controller.CreateContainer(
=======
	containerID, err := controller.CreateContainer(
>>>>>>> master
		container.Config{
			Image: "traefik:1.7-alpine",
			ExposedPorts: nat.PortSet{
				"80/tcp":   struct{}{},
				"8080/tcp": struct{}{},
				"443/tcp":  struct{}{},
			},
			Cmd: cmd,
		},
		container.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
				config.Cfg.Config.GetString("acmePath") + ":/var/acme",
			},
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{HostPort: "80"},
				},
				"8080/tcp": []nat.PortBinding{
					{HostPort: "8080"},
				},
				"443/tcp": []nat.PortBinding{
					{HostPort: "443"},
				},
			},
		},
		network.NetworkingConfig{},
		common.SpaasName("traefik"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
<<<<<<< HEAD
	if err := controller.StartContainer(container.ID); err != nil {
=======
	if err := controller.StartContainer(containerID.ID); err != nil {
>>>>>>> master
		log.Fatal(err.Error())
	}
	log.Println("Traefik is now installed")
}
