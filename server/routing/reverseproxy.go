package routing

import (
	"log"

	client "github.com/fsouza/go-dockerclient"
	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/server/controller"
)

func InitReverseProxy() {
	list, err := controller.ListContainers()
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range list {
		for _, element := range container.Names {
			log.Println(element)
			if element == "/"+common.SpaasName("nginx-proxy") {
				return
			}
		}
	}
	log.Println("Nginx Proxy is not installed but installing now")
	if err := controller.PullImage("jwilder/nginx-proxy", "alpine-0.7.0"); err != nil {
		log.Fatal(err.Error())
	}
	containerID, err := controller.CreateContainer(client.CreateContainerOptions{
		Name: common.SpaasName("nginx-proxy"),
		Config: &client.Config{
			Image: "jwilder/nginx-proxy:alpine-0.7.0",
			ExposedPorts: map[client.Port]struct{}{
				"80/tcp": struct{}{},
			},
		},
		HostConfig: &client.HostConfig{
			Binds: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
			PortBindings: map[client.Port][]client.PortBinding{
				"80/tcp": []client.PortBinding{
					{HostIP: "0.0.0.0", HostPort: "80"},
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
	log.Println("Nginx Proxy is now installed")
}
