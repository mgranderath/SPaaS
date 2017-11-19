package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	scribble "github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver

// InitDB : Initialize the database connection
func InitDB() {
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, "Database")
	db, err = scribble.New(path, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
