package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	"github.com/magrandera/PiaaS/server/docker"
	git "gopkg.in/src-d/go-git.v4"
)

// Application : stores information about the applications
type Application struct {
	Name        string `json:"name"`
	Running     bool   `json:"running"`
	Path        string `json:"path"`
	Repository  string `json:"repo"`
	Type        string `json:"type"`
	Port        string `json:"port"`
	ContainerID string `json:"containerID"`
}

// CreateApplication creates a new application
func CreateApplication(name string) (Application, error) {
	home := GetHomeFolder()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name, "repo")
	// Check whether folder exists
	if val, _ := exists(path); val {
		return Application{}, errors.New("already exists")
	}
	// Create the folders
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return Application{}, err
	}
	// Initlialize the git repository
	if _, err := git.PlainInit(path, true); err != nil {
		return Application{}, err
	}
	// Create the post-receive hook
	err = os.MkdirAll(filepath.Join(path, "hooks"), os.ModePerm)
	if err != nil {
		return Application{}, err
	}
	file, err := os.Create(filepath.Join(path, "hooks", "post-receive"))
	if err != nil {
		return Application{}, err
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\necho \"Starting deploy!\"\n"+
		"curl --request POST 'http://127.0.0.1:5000/api/v1/app/%s/deploy' |"+
		"python -c 'import json,sys;obj=json.load(sys.stdin);print \"Successfully deployed!\";print \"Port: \"+obj[\"port\"]'\n", name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		return Application{}, err
	}
	// Initialize the database record
	app := Application{}
	app.Name = name
	app.Path = filepath.Join(basePath, "Applications", name)
	app.Repository = path
	if err := db.Write("app", name, app); err != nil {
		return Application{}, err
	}
	return app, nil
}

// DeleteApplication deletes existing application
func DeleteApplication(name string) (bool, error) {
	home := GetHomeFolder()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name)
	// Check whether app exists
	if _, err := os.Stat(path); err != nil {
		return false, err
	}
	dock, err := docker.New()
	if err != nil {
		return false, err
	}
	// Remove the container
	if err = dock.RemoveContainer(piName(name)); err != nil {
		return false, err
	}
	// Remove the docker image
	if err = dock.RemoveImage(piName(name)); err != nil {
		return false, err
	}
	// Remove directories
	err = os.RemoveAll(path)
	// Delete app from the database
	if err := db.Delete("app", name); err != nil {
		return false, err
	}
	return true, nil
}

// DeployApplication deploys an application
func DeployApplication(name string) (Application, error) {
	// Create deploy folder
	app, err := GetApplication(name)
	if err != nil {
		return Application{}, err
	}
	path := filepath.Join(app.Path, "deploy")
	if err = os.RemoveAll(path); err != nil {
		return Application{}, err
	}
	// Clone repository into the deploy directory
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		return Application{}, err
	}
	// Parse the Procfile
	dockerfile := Dockerfile{}
	proc := ParseProcfile(path + "/Procfile")
	for _, el := range proc {
		if el.Name == "web" {
			commands := strings.Fields(el.Command)
			dockerfile.Command = commands
		}
	}
	// Detect the language
	if FileExists(filepath.Join(path, "requirements.txt")) {
		app.Type = "python"
	} else if FileExists(filepath.Join(path, "package.json")) {
		app.Type = "nodejs"
	} else if FileExists(filepath.Join(path, "Gemfile")) {
		app.Type = "ruby"
	} else {
		return Application{}, errors.New("no type detected")
	}
	// Search for a free port
	l, _ := net.Listen("tcp", ":0")
	hostport := l.Addr().String()
	_, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return Application{}, err
	}
	dockerfile.Port = port
	app.Port = port
	l.Close()
	// Create Dockerfile
	err = CreateDockerfile(dockerfile, app)
	if err != nil {
		return Application{}, err
	}
	// Tar the deploy directory
	session := sh.NewSession()
	session.SetDir(path + "/")
	_, err = session.Command("tar", "cvf", "../package.tar", ".").Output()
	if err != nil {
		return Application{}, err
	}
	f, err := os.Open(filepath.Join(path, "..", "package.tar"))
	if err != nil {
		return Application{}, err
	}
	defer f.Close()
	// Create connection to Docker API
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	buildResponse, err := dock.BuildImage(f, piName(name))
	if err != nil {
		return Application{}, err
	}
	// Reformat the output
	type Event struct {
		Stream string `json:"stream"`
	}
	d := json.NewDecoder(buildResponse.Body)
	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		event.Stream = strings.TrimSuffix(event.Stream, "\n")
		if strings.Contains(event.Stream, "Step") {
		}
	}
	defer buildResponse.Body.Close()
	// Remove old container
	if err = dock.RemoveContainer(piName(name)); err != nil {
		return Application{}, err
	}
	// Create Container
	createResponse, err := dock.BuildContainer(piName(name), piName(name), name, port)
	if err != nil {
		return Application{}, err
	}
	app.ContainerID = createResponse.ID
	// Start the container
	if err = dock.StartContainer(piName(name)); err != nil {
		return Application{}, err
	}
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		return Application{}, err
	}
	return app, nil
}

// StopApplication stops the container of an application
func StopApplication(name string) (Application, error) {
	app, err := GetApplication(name)
	if err != nil {
		return Application{}, err
	}
	if app.Name == "" {
		return Application{}, errors.New("app does not exist")
	}
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	if err := dock.StopContainer(piName(name)); err != nil {
		return Application{}, err
	}
	app.Running = false
	if err := db.Write("app", name, app); err != nil {
		return Application{}, err
	}
	return app, nil
}

// StartApplication starts an application
func StartApplication(name string) (Application, error) {
	app, err := GetApplication(name)
	if err != nil {
		return Application{}, err
	}
	if app.Name == "" {
		return Application{}, errors.New("app does not exist")
	}
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	if err = dock.StartContainer(piName(name)); err != nil {
		return Application{}, err
	}
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		return Application{}, err
	}
	return app, nil
}

// GetApplication get specific application
func GetApplication(name string) (Application, error) {
	app := Application{}
	if err := db.Read("app", name, &app); err != nil {
		return Application{}, err
	}
	return app, nil
}

// GetApplications get a list of all applications
func GetApplications() ([]Application, error) {
	records, err := db.ReadAll("app")
	if err != nil {
		fmt.Println("Error", err)
		return []Application{}, err
	}
	applications := []Application{}
	for _, f := range records {
		appFound := Application{}
		if err := json.Unmarshal([]byte(f), &appFound); err != nil {
			fmt.Println("Error", err)
			return []Application{}, err
		}
		applications = append(applications, appFound)
	}
	return applications, nil
}
