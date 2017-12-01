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
	home := getHomeFolder()
	executable := getExecutablePath()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name, "repo")
	// Check whether folder exists
	if _, err := os.Stat(path); err == nil {
		printErr(w, "App already exists.")
		return Application{}, err
	}
	// Create the folders
	printNormal(w, "Creating Application '"+name+"'.")
	printNormal(w, "Creating directories")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	printSuccess(w, "Creating directories.")
	// Initlialize the git repository
	printNormal(w, "Initializing git repository.")
	if _, err := git.PlainInit(path, true); err != nil {
		printErr(w, err)
		return Application{}, err
	}
	// Create the post-receive hook
	err = os.MkdirAll(filepath.Join(path, "hooks"), os.ModePerm)
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	file, err := os.Create(filepath.Join(path, "hooks", "post-receive"))
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	defer file.Close()
	fmt.Fprintf(file, "#!/usr/bin/env bash\n%s/PiaaS app deploy %s\n", executable, name)
	// Make the hook executable
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	printSuccess(w, "Initializing git repository.")
	printInfo(w, "Repository path: "+path)
	// Initialize the database record
	app := Application{}
	app.Name = name
	app.Path = filepath.Join(basePath, "Applications", name)
	app.Repository = path
	printNormal(w, "Creating database record.")
	if err := db.Write("app", name, app); err != nil {
		printErr(w, err)
		return Application{}, err
	}
	printSuccess(w, "Creating database record.")
	printSuccess(w, ("Application '" + name + "' successfully created."))
	return app, nil
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(w *os.File, name string) (bool, error) {
	home := getHomeFolder()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name)
	// Check whether app exists
	if _, err := os.Stat(path); err != nil {
		printErr(w, "App does not exist.")
		return false, err
	}
	printNormal(w, ("Removing Application '" + name + "'."))
	printNormal(w, "Deleting Containers.")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		printErr(w, err)
		return false, err
	}
	// Remove the docker container
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Deleting Containers.")
	printNormal(w, "Deleting Images.")
	// Remove the docker image
	_, err = cli.ImageRemove(ctx, name, types.ImageRemoveOptions{Force: true})
	if err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Deleting Images.")
	// Remove directories
	printNormal(w, "Removing Directories")
	if err := os.RemoveAll(path); err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Removing directories.")
	// Delete app from the database
	printNormal(w, "Deleting Database Record.")
	if err := db.Delete("app", name); err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Deleting database record")
	printSuccess(w, ("Application '" + name + "' successfully removed."))
	return true, nil
}

// DeployApplication : Deploys the application
func DeployApplication(w *os.File, name string) {
	// Create deploy folder
	printNormal(w, ("Deploying '" + name + "'."))
	app, err := GetApplication(name)
	if err != nil {
		printErr(w, err)
		return
	}
	path := filepath.Join(app.Path, "deploy")
	if err = os.RemoveAll(path); err != nil {
		printErr(w, err)
		return
	}
	// Clone repository into the deploy directory
	printNormal(w, "Creating seperate directory for deployment.")
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		printErr(w, err)
		return
	}
	printSuccess(w, "Creating seperate directory for deployment.")
	// Parse the Procfile
	dock := Dockerfile{}
	proc := parseProcfile(path + "/Procfile")
	for _, el := range proc {
		if el.Name == "web" {
			commands := strings.Fields(el.Command)
			dock.Command = commands
		}
	}
	// Detect the language
	if fileExists(filepath.Join(path, "requirements.txt")) {
		printInfo(w, "Python was detected")
		app.Type = "python"
	} else if fileExists(filepath.Join(path, "package.json")) {
		printInfo(w, "NodeJs was detected")
		app.Type = "nodejs"
	} else {
		printErr(w, "No type detected.")
		return
	}
	// Search for a free port
	printNormal(w, "Detecting free port")
	l, _ := net.Listen("tcp", ":0")
	hostport := l.Addr().String()
	_, port, err := net.SplitHostPort(hostport)
	if err != nil {
		printErr(w, err)
		return
	}
	dock.Port = port
	app.Port = port
	l.Close()
	printInfo(w, "Port allocated: "+port)
	// Create Dockerfile
	printNormal(w, "Creating Dockerfile")
	err = CreateDockerfile(dock, app)
	if err != nil {
		printErr(w, err)
		return
	}
	printSuccess(w, "Creating Dockerfile")
	printNormal(w, "Creating Docker image")
	// Tar the deploy directory
	session := sh.NewSession()
	session.Stdout = w
	session.Stderr = w
	session.SetDir(path + "/")
	_, err = session.Command("tar", "cvf", "../package.tar", ".").Output()
	if err != nil {
		printErr(w, err)
		return
	}
	f, err := os.Open(filepath.Join(path, "..", "package.tar"))
	if err != nil {
		printErr(w, err)
	}
	defer f.Close()
	// Create connection to Docker API
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		printErr(w, err)
		return
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
		printErr(w, err)
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
			printInfo(w, event.Stream)
		}
	}
	defer imageBuildResponse.Body.Close()
	printSuccess(w, "Creating Docker image")
	printNormal(w, "Creating Container")
	// Remove old container
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true})
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		printErr(w, err)
	}
	// Create Container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: name,
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
		printErr(w, err)
	}
	app.ContainerID = resp.ID
	printSuccess(w, "Creating Container")
	printNormal(w, "Starting Container")
	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		printErr(w, err)
	}
	printSuccess(w, "Starting Container")
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		printErr(w, err)
		return
	}
}

// StopApplication : Stops the container of a application
func StopApplication(w *os.File, name string) error {
	printNormal(w, "Stopping Application '"+name+"'.")
	app, err := GetApplication(name)
	if err != nil {
		printErr(w, err)
		return err
	}
	if app.Name == "" {
		printErr(w, name+" does not exist!")
		return errors.New("app does not exist")
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		printErr(w, err)
		return err
	}
	err = cli.ContainerStop(ctx, name, nil)
	if err != nil {
		printErr(w, err)
		return err
	}
	printSuccess(w, "Stopping Application '"+name+"'.")
	app.Running = false
	if err := db.Write("app", name, app); err != nil {
		printErr(w, err)
		return err
	}
	return nil
}

// StartApplication : starts the application
func StartApplication(w *os.File, name string) error {
	printNormal(w, "Starting Application '"+name+"'.")
	app, err := GetApplication(name)
	if err != nil {
		printErr(w, err)
		return err
	}
	if app.Name == "" {
		printErr(w, name+" does not exist!")
		return errors.New("app does not exist")
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		printErr(w, err)
		return err
	}
	err = cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
	if err != nil {
		printErr(w, err)
		return err
	}
	printSuccess(w, "Starting Application '"+name+"'.")
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		printErr(w, err)
		return err
	}
	return nil
}

// LogApplication : show the log oft the application
func LogApplication(w *os.File, name string, tail bool) {
	printNormal(w, "Logs of '"+name+"'.")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		printErr(w, err)
		return
	}
	out, err := cli.ContainerLogs(ctx, name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     tail})
	if err != nil {
		printErr(w, err)
		return
	}
	io.Copy(os.Stdout, out)
}

// GetApplication : Get specific application
func GetApplication(name string) (Application, error) {
	app := Application{}
	if err := db.Read("app", name, &app); err != nil {
		printErr(os.Stdout, err)
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
