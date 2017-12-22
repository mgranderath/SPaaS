package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
	executable := getExecutablePath()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name, "repo")
	// Check whether folder exists
	if _, err := os.Stat(path); err == nil {
		PrintErr(w, "App already exists.")
		return Application{}, err
	}
	// Create the folders
	PrintNormal(w, "Creating Application '"+name+"'.")
	PrintNormal(w, "Creating directories")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Creating directories.")
	// Initlialize the git repository
	PrintNormal(w, "Initializing git repository.")
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
	fmt.Fprintf(file, "#!/usr/bin/env bash\n%s/PiaaS app deploy %s\n", executable, name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Initializing git repository.")
	PrintInfo(w, "Repository path: "+path)
	// Initialize the database record
	app := Application{}
	app.Name = name
	app.Path = filepath.Join(basePath, "Applications", name)
	app.Repository = path
	PrintNormal(w, "Creating database record.")
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Creating database record.")
	PrintSuccess(w, ("Application '" + name + "' successfully created."))
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
	PrintNormal(w, ("Removing Application '" + name + "'."))
	PrintNormal(w, "Deleting Containers.")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		PrintErr(w, err)
		return false, err
	}
	// Remove the container
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		PrintErr(w, err)
	}
	PrintSuccess(w, "Deleting Containers.")
	PrintNormal(w, "Deleting Images.")
	// Remove the docker image
	_, err = cli.ImageRemove(ctx, name, types.ImageRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such image") {
		PrintErr(w, err)
		return false, err
	}
	PrintSuccess(w, "Deleting Images.")
	// Remove directories
	PrintNormal(w, "Removing Directories")
	if err := os.RemoveAll(path); err != nil {
		PrintErr(w, err)
		return false, err
	}
	PrintSuccess(w, "Removing directories.")
	// Delete app from the database
	PrintNormal(w, "Deleting Database Record.")
	if err := db.Delete("app", name); err != nil {
		PrintErr(w, err)
		return false, err
	}
	PrintSuccess(w, "Deleting database record")
	PrintSuccess(w, ("Application '" + name + "' successfully removed."))
	return true, nil
}

// DeployApplication : Deploys the application
func DeployApplication(w *os.File, name string) (Application, error) {
	// Create deploy folder
	PrintNormal(w, ("Deploying '" + name + "'."))
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
	PrintNormal(w, "Creating seperate directory for deployment.")
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Creating seperate directory for deployment.")
	// Parse the Procfile
	dock := Dockerfile{}
	proc := ParseProcfile(path + "/Procfile")
	for _, el := range proc {
		if el.Name == "web" {
			commands := strings.Fields(el.Command)
			dock.Command = commands
		}
	}
	// Detect the language
	if fileExists(filepath.Join(path, "requirements.txt")) {
		PrintInfo(w, "Python was detected")
		app.Type = "python"
	} else if fileExists(filepath.Join(path, "package.json")) {
		PrintInfo(w, "NodeJs was detected")
		app.Type = "nodejs"
	} else if fileExists(filepath.Join(path, "Gemfile")) {
		PrintInfo(w, "Ruby was detected")
		app.Type = "ruby"
	} else {
		PrintErr(w, "No type detected.")
		return Application{}, errors.New("no type detected")
	}
	// Search for a free port
	PrintNormal(w, "Detecting free port")
	l, _ := net.Listen("tcp", ":0")
	hostport := l.Addr().String()
	_, port, err := net.SplitHostPort(hostport)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	dock.Port = port
	app.Port = port
	l.Close()
	PrintInfo(w, "Port allocated: "+port)
	// Create Dockerfile
	PrintNormal(w, "Creating Dockerfile")
	err = CreateDockerfile(dock, app)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Creating Dockerfile")
	PrintNormal(w, "Creating Docker image")
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
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Build the image
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		f,
		types.ImageBuildOptions{
			Tags:       []string{name},
			Dockerfile: "Dockerfile",
			Remove:     true,
			NoCache:    true})
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	// Reformat the output
	type Event struct {
		Stream string `json:"stream"`
	}
	d := json.NewDecoder(imageBuildResponse.Body)
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
			PrintInfo(w, event.Stream)
		}
	}
	defer imageBuildResponse.Body.Close()
	PrintSuccess(w, "Creating Docker image")
	PrintNormal(w, "Creating Container")
	// Remove old container
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		PrintErr(w, err)
		return Application{}, err
	}
	// Create Container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: name,
		Env:   []string{"VIRTUAL_HOST=" + name + ".granderath.tech"},
		ExposedPorts: nat.PortSet{
			"5000/tcp": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"5000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		}}, nil, name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	app.ContainerID = resp.ID
	PrintSuccess(w, "Creating Container")
	PrintNormal(w, "Starting Container")
	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Starting Container")
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// StopApplication : Stops the container of a application
func StopApplication(w *os.File, name string) (Application, error) {
	PrintNormal(w, "Stopping Application '"+name+"'.")
	app, err := GetApplication(name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	if app.Name == "" {
		PrintErr(w, name+" does not exist!")
		return Application{}, errors.New("app does not exist")
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	err = cli.ContainerStop(ctx, name, nil)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Stopping Application '"+name+"'.")
	app.Running = false
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// StartApplication : starts the application
func StartApplication(w *os.File, name string) (Application, error) {
	PrintNormal(w, "Starting Application '"+name+"'.")
	app, err := GetApplication(name)
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	if app.Name == "" {
		PrintErr(w, name+" does not exist!")
		return Application{}, errors.New("app does not exist")
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	err = cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
	if err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	PrintSuccess(w, "Starting Application '"+name+"'.")
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		PrintErr(w, err)
		return Application{}, err
	}
	return app, nil
}

// LogApplication : show the log oft the application
func LogApplication(w *os.File, name string, tail bool) {
	PrintNormal(w, "Logs of '"+name+"'.")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		PrintErr(w, err)
		return
	}
	out, err := cli.ContainerLogs(ctx, name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     tail})
	if err != nil {
		PrintErr(w, err)
		return
	}
	io.Copy(os.Stdout, out)
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
