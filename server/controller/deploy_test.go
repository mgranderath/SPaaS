package controller

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	git "gopkg.in/src-d/go-git.v4"
)

func TestItDeploysApp(t *testing.T) {
	messages := make(chan Application)
	name := "test"
	config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	InitDocker()
	config.Cfg.Config.Set("HOST_CONFIG_FOLDER", filepath.Join(common.HomeDir(), ".spaas"))
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	// Setting up
	err := os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		testRemoveApp(appPath, t)
		t.Fatal(err.Error())
	}
	_, err = git.PlainClone(repoPath, true, &git.CloneOptions{
		URL: "https://github.com/magrandera/SPaaS-node-js-example.git",
	})
	if err != nil {
		testRemoveApp(appPath, t)
		t.Fatal(err.Error())
	}
	go deploy(name, messages)
	for elem := range messages {
		if elem.Type == "error" {
			testRemoveApp(appPath, t)
			t.Fatal("Gave back an error: " + err.Error())
		}
	}
	container, err := dock.Cli.ContainerInspect(dock.Ctx, common.SpaasName(name))
	if err != nil {
		testRemoveApp(appPath, t)
		testRemoveContainer(common.SpaasName(name), t)
		testRemoveImage(common.SpaasName(name), t)
		t.Fatal("Gave back an error: " + err.Error())
	}
	if !container.State.Running {
		testRemoveApp(appPath, t)
		testRemoveContainer(common.SpaasName(name), t)
		testRemoveImage(common.SpaasName(name), t)
		t.Fatal("Gave back an error: " + err.Error())
	}
	testRemoveApp(appPath, t)
	testRemoveContainer(common.SpaasName(name), t)
	testRemoveImage(common.SpaasName(name), t)
}
