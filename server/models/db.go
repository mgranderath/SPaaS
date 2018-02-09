package models

import (
	"fmt"
	"path/filepath"

	scribble "github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver

// InitDB initializes the database connection
func InitDB() {
	var err error
	dir := filepath.Join(GetHomeFolder(), "PiaaS-Data")
	path := filepath.Join(dir, "Database")
	db, err = scribble.New(path, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
