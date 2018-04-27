package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/magrandera/PiaaS/server/config"
	"github.com/magrandera/PiaaS/server/docker"

	"github.com/magrandera/PiaaS/server/models"

	"github.com/gorilla/mux"
	"github.com/magrandera/PiaaS/server/routing"
)

func initServer() error {
	configPath := filepath.Join(models.GetHomeFolder(), ".config/piaas/config")
	if !models.FileExists(configPath) {
		configFile := config.Configuration{}
		configFile.Nginx = false
		configFile.Secret = models.RandString(32)
		err := config.WriteConfig(configPath, configFile)
		if err != nil {
			return err
		}
	}
	config, err := config.ReadConfig(configPath)
	if err != nil {
		return err
	}
	if !config.Nginx {
		return nil
	}
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

var mySigningKey = []byte("secret")

// StartServer starts the PiaaS server
func StartServer() {
	router := mux.NewRouter()
	routing.SetupRouting(router)
	models.InitDB()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user"] = "magrandera"

	tokenString, _ := token.SignedString(mySigningKey)
	fmt.Println(tokenString)

	err := initServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(http.ListenAndServe(":5000", router))
}
