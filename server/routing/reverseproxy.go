package routing

import (
	"log"

	client "github.com/fsouza/go-dockerclient"
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
	if err := controller.PullImage("traefik", "1.7-alpine"); err != nil {
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
	containerID, err := controller.CreateContainer(client.CreateContainerOptions{
		Name: common.SpaasName("traefik"),
		Config: &client.Config{
			Image: "traefik:1.7-alpine",
			ExposedPorts: map[client.Port]struct{}{
				"80/tcp":   struct{}{},
				"8080/tcp": struct{}{},
				"443/tcp":  struct{}{},
			},
			Cmd: cmd,
		},
		HostConfig: &client.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
				config.Cfg.Config.GetString("acmePath") + ":/var/acme",
			},
			PortBindings: map[client.Port][]client.PortBinding{
				"80/tcp": []client.PortBinding{
					{HostPort: "80"},
				},
				"8080/tcp": []client.PortBinding{
					{HostPort: "8080"},
				},
				"443/tcp": []client.PortBinding{
					{HostPort: "443"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := controller.StartContainer(containerID.ID); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Traefik is now installed")
}
