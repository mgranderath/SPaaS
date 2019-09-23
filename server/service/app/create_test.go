package app

import (
	"github.com/mgranderath/SPaaS/server/model"
	"os"
	"path/filepath"
	"testing"

	"docker.io/go-docker/api/types"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
)

func testRemoveApp(appPath string, t *testing.T) {
	if err := os.RemoveAll(appPath); err != nil {
		t.Fatal(err.Error())
	}
}

func testRemoveContainer(name string, t *testing.T) {
	err := dock.Cli.ContainerRemove(dock.Ctx, name, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func testRemoveImage(name string, t *testing.T) {
	_, err := dock.Cli.ImageRemove(dock.Ctx, name, types.ImageRemoveOptions{
		Force: true,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestItCreatesNewApp(t *testing.T) {
	messages := make(chan model.Status)
	name := "test"
	config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	config.ConfigStore.Config.Set("HOST_CONFIG_FOLDER", filepath.Join(common.HomeDir(), ".spaas"))
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	go create(name, messages)
	for elem := range messages {
		if elem.Type == "error" {
			testRemoveApp(appPath, t)
			t.Fatal("Gave back an error: " + elem.Message)
		}
	}
	if !common.Exists(repoPath) {
		t.Fatal("Did not create the repository")
	}
	testRemoveApp(appPath, t)
}
