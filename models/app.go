package models

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	git "gopkg.in/src-d/go-git.v4"
)

// Application : stores information about the applications
type Application struct {
	Name       string `json:"name"`
	Running    bool   `json:"running"`
	Path       string `json:"path"`
	Repository string `json:"repo"`
	Type       string `json:"type"`
}

// CreateApplication : Creates a new Application
func CreateApplication(w *os.File, name string) (Application, error) {
	home := getHomeFolder()
	executable := getExecutablePath()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name, "repo")

	if _, err := os.Stat(path); err == nil {
		printErr(w, "App already exists.")
		return Application{}, err
	}
	printNormal(w, "Creating Application '"+name+"'.")
	printNormal(w, "Creating directories")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	printSuccess(w, "Creating directories.")

	printNormal(w, "Initializing git repository.")
	if _, err := git.PlainInit(path, true); err != nil {
		printErr(w, err)
		return Application{}, err
	}
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
	err = os.Chmod(filepath.Join(path, "hooks", "post-receive"), 0755)
	if err != nil {
		printErr(w, err)
		return Application{}, err
	}
	printSuccess(w, "Initializing git repository.")
	printInfo(w, "Repository path: "+path)

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
	session := sh.NewSession()
	basePath := filepath.Join(home, "PiaaS-Data")
	path := filepath.Join(basePath, "Applications", name)
	if _, err := os.Stat(path); err != nil {
		printErr(w, "App does not exist.")
		return false, err
	}
	printNormal(w, ("Removing Application '" + name + "'."))
	printNormal(w, "Deleting Containers.")
	_, err := session.Command("docker", "rm", "--force", name).Output()
	if err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Deleting Images.")
	printNormal(w, "Deleting Images.")
	_, err = session.Command("docker", "rmi", name).Output()
	if err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Deleting Images.")
	printNormal(w, "Removing Directories")
	if err := os.RemoveAll(path); err != nil {
		printErr(w, err)
		return false, err
	}
	printSuccess(w, "Removing directories.")
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
	printNormal(w, "Creating seperate directory for deployment.")
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: app.Repository,
	})
	if err != nil {
		printErr(w, err)
		return
	}
	printSuccess(w, "Creating seperate directory for deployment.")
	dock := Dockerfile{}
	proc := parseProcfile(path + "/Procfile")
	for _, el := range proc {
		if el.Name == "web" {
			commands := strings.Fields(el.Command)
			dock.Command = commands
		}
	}
	if fileExists(filepath.Join(path, "requirements.txt")) {
		printInfo(w, "Python was detected")
		app.Type = "python"
		if err := db.Write("app", name, app); err != nil {
			printErr(w, err)
			return
		}
	} else {
		printErr(w, "No type detected.")
		return
	}
	l, _ := net.Listen("tcp", ":0")
	hostport := l.Addr().String()
	_, port, err := net.SplitHostPort(hostport)
	if err != nil {
		printErr(w, err)
		return
	}
	dock.Port = port
	l.Close()
	printInfo(w, port)
	err = CreateDockerfile(dock, app)
	if err != nil {
		printErr(w, err)
		return
	}
	session := sh.NewSession()
	session.Stdout = w
	session.Stderr = w
	session.SetDir(path + "/")
	_, err = session.Command("docker", "build", "-t", name, ".").Output()
	if err != nil {
		printErr(w, err)
		return
	}
	_, err = session.Command("docker", "run", "-d", "-p", port+":5000", "--name", name, name).Output()
	if err != nil {
		printErr(w, err)
		return
	}
	app.Running = true
	if err := db.Write("app", name, app); err != nil {
		printErr(w, err)
		return
	}
}

// StopApplication : Stops the container of a application
func StopApplication(w *os.File, name string) error {
	printNormal(w, "Stopping Application '"+name+"'.")
	session := sh.NewSession()
	_, err := session.Command("docker", "stop", name).Output()
	if err != nil {
		printErr(w, err)
		return err
	}
	printSuccess(w, "Stopping Application '"+name+"'.")
	return nil
}

// StartApplication : starts the application
func StartApplication(w *os.File, name string) error {
	printNormal(w, "Starting Application '"+name+"'.")
	session := sh.NewSession()
	_, err := session.Command("docker", "start", name).Output()
	if err != nil {
		printErr(w, err)
		return err
	}
	printSuccess(w, "Starting Application '"+name+"'.")
	return nil
}

// LogApplication : show the log oft the application
func LogApplication(w *os.File, name string, tail bool) {
	printNormal(w, "Logs of '"+name+"'.")
	session := sh.NewSession()
	session.Stdout = w
	session.Stderr = w
	if tail {
		session.Command("docker", "logs", "--follow", name).Run()
	} else {
		session.Command("docker", "logs", name).Run()
	}
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
