package services

import (
	"docker.io/go-docker/api/types"
	"github.com/mgranderath/SPaaS/common"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (a *AppService) GetApplicationStats(name string) (types.ContainerJSON, error) {
	return a.DockerClient.InspectContainer(common.SpaasName(name))
}

func (a *AppService) GetApplications() ([]string, error) {
	appPath := filepath.Join(basePath, "applications")
	files, err := ioutil.ReadDir(appPath)
	if err != nil {
		return []string{}, err
	}
	return filterForDir(files), nil
}

func filterForDir(files []os.FileInfo) []string {
	filtered := make([]string, 0)
	for _, item := range files {
		if item.IsDir() {
			filtered = append(filtered, item.Name())
		}
	}
	return filtered
}
