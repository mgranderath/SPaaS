package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type application struct {
	ID   string
	Name string
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
	os.MkdirAll(path, os.ModePerm)
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
}

// GetApplications : Get a list of all applications
func GetApplications(w http.ResponseWriter, r *http.Request) {

}

// GetApplication : Get specific application
func GetApplication(w http.ResponseWriter, r *http.Request) {

}
