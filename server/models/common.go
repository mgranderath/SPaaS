package models

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	color "github.com/logrusorgru/aurora"
)

func fileExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// GetHomeFolder : get filepath to home folder of user
func GetHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func getExecutablePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// PrintErr : prints an error
func PrintErr(w *os.File, args interface{}) {
	switch args.(type) {
	case string:
		fmt.Fprintln(w, color.Red("-----> ERROR: "+args.(string)))
		break
	case error:
		fmt.Fprintln(w, color.Red("-----> ERROR: "+args.(error).Error()))
		break
	}
}

// PrintSuccess : prints a success message
func PrintSuccess(w *os.File, message string) {
	fmt.Fprintln(w, color.Green("-----> Success: "+message))
}

// PrintNormal : prints a normal message
func PrintNormal(w *os.File, message string) {
	fmt.Fprintln(w, ("-----> Task: " + message))
}

// PrintInfo : prints a info message
func PrintInfo(w *os.File, message string) {
	fmt.Fprintln(w, color.Brown("-----> Info: "+message))
}
