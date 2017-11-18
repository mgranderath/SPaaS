package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	git "gopkg.in/src-d/go-git.v4"
)

type application struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
}

// CreateApplication : Creates a new Application
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	appPath := filepath.Join("Applications", params["name"])
	path := filepath.Join(dir, appPath)
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			fmt.Fprint(w, "Already exists!")
			return
		}
	}
	os.MkdirAll(path, os.ModePerm)
	git.PlainInit(path, true)
	fmt.Fprint(w, path)
	app := application{}
	if err := db.Write("app", params["name"], app); err != nil {
		fmt.Println("Error", err)
	}
}

// DeleteApplication : Deletes existing Application
func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	appPath := filepath.Join("Applications", params["name"])
	path := filepath.Join(dir, appPath)
	os.RemoveAll(path)
	if err := db.Delete("app", params["name"]); err != nil {
		fmt.Println("Error", err)
	}
}

// GetApplications : Get a list of all applications
func GetApplications(w http.ResponseWriter, r *http.Request) {
	records, err := db.ReadAll("app")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Fprint(w, records)
}

// UpdateApplication : Update the state of an application
func UpdateApplication(w http.ResponseWriter, r *http.Request) {

}

// GetApplication : Get specific application
func GetApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app := application{}
	if err := db.Read("app", params["name"], &app); err != nil {
		fmt.Println("Error", err)
	}
	b, err := json.Marshal(app)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(b[:]))
}
