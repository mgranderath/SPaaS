package models

import (
	"encoding/json"
	"fmt"
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
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return Application{}, err
	}
	appPath := filepath.Join("Applications", name)
	path := filepath.Join(dir, appPath)
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintln(w, "----->ERROR: App already exists.")
		return Application{}, err
	}
	fmt.Fprintf(w, "----->Creating Application '%s'.\n", name)
	fmt.Fprintln(w, "----->Creating directories.")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success creating directories.")
	fmt.Fprintln(w, "----->Initializing git repository.")
	if _, err := git.PlainInit(path, true); err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success initializing git repository.")
	fmt.Fprintf(w, "----->Repository path: %s\n", path)
	app := Application{}
	app.Name = name
	app.Path = path
	fmt.Fprintln(w, "----->Creating database record.")
	if err := db.Write("app", name, app); err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return Application{}, err
	}
	fmt.Fprintln(w, "----->Success creating database record.")
	fmt.Fprintf(w, "----->Application '%s' successfully created.\n", name)
	return app, nil
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(w *os.File, name string) (bool, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return false, err
	}
	appPath := filepath.Join("Applications", name)
	path := filepath.Join(dir, appPath)
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintln(w, "----->ERROR: App does not exist.")
		return false, err
	}
	fmt.Fprintf(w, "----->Removing Application '%s'.\n", name)
	fmt.Fprintln(w, "----->Removing Directories.")
	if err = os.RemoveAll(path); err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return false, err
	}
	fmt.Fprintln(w, "----->Success removing directories.")
	fmt.Fprintln(w, "----->Deleting Database Record.")
	if err := db.Delete("app", name); err != nil {
		fmt.Fprint(w, "----->")
		fmt.Fprintln(w, err)
		return false, err
	}
	fmt.Fprintln(w, "----->Success deleting database record.")
	fmt.Fprintf(w, "----->Application '%s' successfully removed.\n", name)
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
	//TODO: Implement update function
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
