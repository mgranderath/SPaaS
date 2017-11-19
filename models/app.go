package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"
)

// Application : stores information about the applications
type Application struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Path    string `json:"path"`
}

// CreateApplication : Creates a new Application
func CreateApplication(w *os.File, name string) (Application, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return Application{}, err
	}
	appPath := filepath.Join("Applications", name)
	path := filepath.Join(dir, appPath)
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			fmt.Println(err)
			return Application{}, err
		}
	}
	fmt.Fprintln(w, "-----> Creating Directories.")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return Application{}, err
	}
	fmt.Fprintln(w, "-----> Initializing git repository.")
	if _, err := git.PlainInit(path, true); err != nil {
		fmt.Println(err)
		return Application{}, err
	}
	app := Application{}
	app.Name = name
	app.Path = path
	fmt.Fprintln(w, "-----> Creating Database Record.")
	if err := db.Write("app", name, app); err != nil {
		fmt.Println("Error", err)
		return Application{}, err
	}
	return app, nil
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(name string) (bool, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	appPath := filepath.Join("Applications", name)
	path := filepath.Join(dir, appPath)
	if err = os.RemoveAll(path); err != nil {
		fmt.Println("Error", err)
		return false, err
	}
	if err := db.Delete("app", name); err != nil {
		fmt.Println("Error", err)
		return false, err
	}
	return true, nil
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

// UpdateApplication : Update the state of an application
func UpdateApplication(w http.ResponseWriter, r *http.Request) {

}

// GetApplication : Get specific application
func GetApplication(name string) (Application, error) {
	app := Application{}
	if err := db.Read("app", name, &app); err != nil {
		fmt.Println("Error", err)
		return Application{}, err
	}
	return app, nil
}
