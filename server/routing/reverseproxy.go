package routing

import (
	"github.com/mgranderath/SPaaS/server/di"
	"log"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/mgranderath/SPaaS/common"
)

// InitReverseProxy initializes the reverse proxy
func InitReverseProxy(p di.Provider) {
	list, err := p.GetDockerClient().ListContainers()
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
	if err := p.GetDockerClient().PullImage("traefik:1.7-alpine"); err != nil {
		log.Fatal(err.Error())
	}
	cmd := []string{
		"--docker", "--docker.watch",
		"--defaultEntryPoints=http",
		"--entryPoints=Name:http Address::80 Compress:off",
		"--docker.domain=" + p.GetConfigRepository().Config.GetString("domain"),
		"--debug",
		"--logLevel=DEBUG",
	}
	if p.GetConfigRepository().Config.GetBool("letsencrypt") {
		letsencrypt := []string{
			"--acme",
			"--acme.email=" + p.GetConfigRepository().Config.GetString("letsencryptEmail"),
			"--acme.storage=/var/acme/acme.json",
			"--acme.httpchallenge.entrypoint=http",
			"--acme.entrypoint=https",
			"--acme.onhostrule=true",
			"--accesslogsfile=/var/acme/access.log",
			"--entryPoints=Name:https Address::443 TLS Compress:off",
			"--entryPoints=Name:http Address::80 Redirect.EntryPoint:https Compress:off",
			"--defaultEntryPoints=https,http",
		}
		cmd = append(cmd, letsencrypt...)
	}
	containerID, err := p.GetDockerClient().CreateContainer(
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
				p.GetConfigRepository().Config.GetString("acmePath") + ":/var/acme",
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
	if err := p.GetDockerClient().StartContainer(containerID.ID); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Traefik is now installed")
}
