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

// CreateApplication : Creates a new Application
func CreateApplication(w *os.File, name string) (Application, error) {
	home := GetHomeFolder()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name, "repo")
	// Check whether folder exists
	if _, err := os.Stat(path); err == nil {
		PrintErr(w, "App already exists.")
		return Application{}, err
	}
	// Create the folders
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Initlialize the git repository
	if _, err := git.PlainInit(path, true); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Create the post-receive hook
	err = os.MkdirAll(filepath.Join(path, "hooks"), os.ModePerm)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	file, err := os.Create(filepath.Join(path, "hooks", "post-receive"))
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\ncurl --request POST 'http://127.0.0.1:5000/api/v1/app/%s/deploy'\n", name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Initialize the database record
	app := Application{}
	app.Name = name
	app.Path = filepath.Join(basePath, "Applications", name)
	app.Repository = path
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(w *os.File, name string) (bool, error) {
	home := GetHomeFolder()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name)
	// Check whether app exists
	if _, err := os.Stat(path); err != nil {
		PrintErr(w, "App does not exist.")
		return false, err
	}
	dock, err := docker.New()
	if err != nil {
		return false, err
	}
	// Remove the container
	if err = dock.RemoveContainer(name); err != nil {
		return false, err
	}
	// Remove the docker image
	if err = dock.RemoveImage(name); err != nil {
		return false, err
	}
	// Remove directories
	if err := os.RemoveAll(path); err != nil {
		PrintErr(w, err)
		return false, err
	}
	// Delete app from the database
	if err := db.Delete("app", name); err != nil {
		PrintErr(w, err)
		return false, err
	}
	return true, nil
}

// DeployApplication : Deploys the application
func DeployApplication(w *os.File, name string) (Application, error) {
	// Create deploy folder
	app, err := GetApplication(name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	path := filepath.Join(app.Path, "deploy")
	if err = os.RemoveAll(path); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Clone repository into the deploy directory
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		PrintErr(w, err)
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
	if fileExists(filepath.Join(path, "requirements.txt")) {
		app.Type = "python"
	} else if fileExists(filepath.Join(path, "package.json")) {
		app.Type = "nodejs"
	} else if fileExists(filepath.Join(path, "Gemfile")) {
		app.Type = "ruby"
	} else {
		PrintErr(w, "No type detected.")
		return Application{}, errors.New("no type detected")
	}
	// Search for a free port
	l, _ := net.Listen("tcp", ":0")
	hostport := l.Addr().String()
	_, port, err := net.SplitHostPort(hostport)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	dockerfile.Port = port
	app.Port = port
	l.Close()
	// Create Dockerfile
	err = CreateDockerfile(dockerfile, app)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Tar the deploy directory
	session := sh.NewSession()
	session.Stdout = w
	session.Stderr = w
	session.SetDir(path + "/")
	_, err = session.Command("tar", "cvf", "../package.tar", ".").Output()
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	f, err := os.Open(filepath.Join(path, "..", "package.tar"))
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	defer f.Close()
	// Create connection to Docker API
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	buildRespoonse, err := dock.BuildImage(f, name)
	if err != nil {
		return Application{}, err
	}
	// Reformat the output
	type Event struct {
		Stream string `json:"stream"`
	}
	d := json.NewDecoder(buildRespoonse.Body)
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
	defer buildRespoonse.Body.Close()
	// Remove old container
	if err = dock.RemoveContainer(name); err != nil {
		return Application{}, err
	}
	// Create Container
	createResponse, err := dock.BuildContainer(name, port)
	if err != nil {
		return Application{}, err
	}
	app.ContainerID = createResponse.ID
	// Start the container
	if err = dock.StartContainer(name); err != nil {
		return Application{}, err
	}
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// StopApplication : Stops the container of a application
func StopApplication(w *os.File, name string) (Application, error) {
	app, err := GetApplication(name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	if app.Name == "" {
		PrintErr(w, name+" does not exist!")
		return Application{}, errors.New("app does not exist")
	}
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	if err := dock.StopContainer(name); err != nil {
		return Application{}, err
	}
	app.Running = false
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// StartApplication : starts the application
func StartApplication(w *os.File, name string) (Application, error) {
	app, err := GetApplication(name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	if app.Name == "" {
		PrintErr(w, name+" does not exist!")
		return Application{}, errors.New("app does not exist")
	}
	dock, err := docker.New()
	if err != nil {
		return Application{}, err
	}
	if err = dock.StartContainer(name); err != nil {
		return Application{}, err
	}
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// GetApplication : Get specific application
func GetApplication(name string) (Application, error) {
	app := Application{}
	if err := db.Read("app", name, &app); err != nil {
		PrintErr(os.Stdout, err)
		return Application{}, err
	}
	return app, nil
}

// GetApplications : Get a list of all applications
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
