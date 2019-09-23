package model

import (
	"errors"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/config"
	buildpackdocker "github.com/mgranderath/SPaaS/server/buildpack/docker"
	buildpacknodejs "github.com/mgranderath/SPaaS/server/buildpack/nodejs"
	buildpackpython "github.com/mgranderath/SPaaS/server/buildpack/python"
	buildpackruby "github.com/mgranderath/SPaaS/server/buildpack/ruby"
	"os"
	"path/filepath"
	"strings"
)

type Application struct {
	Name           string
	Path           string
	DeployPath     string
	RepositoryPath string
	Type           ApplicationType
	Command        []string
}

var TypeToFile = map[ApplicationType]string{
	Python: "requirements.txt",
	Node:   "package.json",
	Ruby:   "Gemfile",
	Docker: "Dockerfile",
}

var TypeToBuild = map[ApplicationType]func(string, []string) error{
	Python: buildpackpython.Build,
	Node:   buildpacknodejs.Build,
	Ruby:   buildpackruby.Build,
	Docker: buildpackdocker.Build,
}

func NewApplication(name string) *Application {
	basePath := filepath.Join(common.HomeDir(), ".spaas")
	Path := filepath.Join(basePath, "applications", name)
	DeployPath := filepath.Join(Path, "deploy")
	RepositoryPath := filepath.Join(Path, "repo")
	return &Application{
		Name:           name,
		Path:           Path,
		DeployPath:     DeployPath,
		RepositoryPath: RepositoryPath,
		Type:           Undefined,
		Command:        nil,
	}
}

func (app *Application) Exists() bool {
	return common.Exists(app.Path)
}

func (app *Application) DetectType() ApplicationType {
	for appType, file := range TypeToFile {
		if common.Exists(filepath.Join(app.DeployPath, file)) {
			app.Type = appType
			return appType
		}
	}
	return app.Type
}

func (app *Application) Build() error {
	if app.Type == Undefined {
		return errors.New("undefined application type")
	}
	if app.Command == nil || (app.Command != nil && len(app.Command) == 0) {
		return errors.New("run command is undefined")
	}
	err := TypeToBuild[app.Type](app.Path, app.Command)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) ResetDeployDir() error {
	if err := os.RemoveAll(app.DeployPath); err != nil {
		return err
	}
	err := os.MkdirAll(app.DeployPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) DetectStartCommand() error {
	v, err := config.ReadConfig(filepath.Join(app.DeployPath, "spaas.json"), map[string]interface{}{})
	if err != nil {
		return err
	}
	if !v.InConfig("start") {
		return errors.New("No 'start' in spaas.json in project")
	}
	app.Command = strings.Fields(v.GetString("start"))
	return nil
}

func (app Application) Setup() error {
	return nil
}
