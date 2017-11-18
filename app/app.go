package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	scribble "github.com/nanobox-io/golang-scribble"
	git "gopkg.in/src-d/go-git.v4"
)

type application struct {
	Name    string
	Running bool
}

var db *scribble.Driver

// InitDB : Initialize the database connection
func InitDB() {
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	appPath := filepath.Join("Applications", "db")
	path := filepath.Join(dir, appPath)
	db, err = scribble.New(path, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
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
	git.PlainInit(path, true)
	fmt.Fprint(w, path)
	fish := application{}
	if err := db.Write("fish", "onefish", fish); err != nil {
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
}

// GetApplications : Get a list of all applications
func GetApplications(w http.ResponseWriter, r *http.Request) {

}

// UpdateApplication : Update the state of an application
func UpdateApplication(w http.ResponseWriter, r *http.Request) {

}

// GetApplication : Get specific application
func GetApplication(w http.ResponseWriter, r *http.Request) {

}
