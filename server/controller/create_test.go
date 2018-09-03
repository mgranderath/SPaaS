package controller

import (
	"path/filepath"
	"testing"

	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
)

func TestItCreatesNewApp(t *testing.T) {
	messages := make(chan Application)
	name := "test"
	config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	config.Cfg.Config.Set("HOST_CONFIG_FOLDER", filepath.Join(common.HomeDir(), ".spaas"))
	appPath := filepath.Join(basePath, "applications", name)
	repoPath := filepath.Join(appPath, "repo")
	go create(name, messages)
	for elem := range messages {
		if elem.Type == "error" {
			t.Fatal("Gave back an error")
		}
	}
	if !common.Exists(repoPath) {
		t.Fatal("Did not create the repository")
	}
}
