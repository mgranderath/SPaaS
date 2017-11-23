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

func getHomeFolder() string {
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

func printErr(w *os.File, args interface{}) {
	switch args.(type) {
	case string:
		fmt.Fprintln(w, color.Red("-----> ERROR: "+args.(string)))
		break
	case error:
		fmt.Fprintln(w, color.Red("-----> ERROR: "+args.(error).Error()))
		break
	}
}

func printSuccess(w *os.File, message string) {
	fmt.Fprintln(w, color.Green("-----> Success: "+message))
}

func printNormal(w *os.File, message string) {
	fmt.Fprintln(w, ("-----> Task: " + message))
}

func printInfo(w *os.File, message string) {
	fmt.Fprintln(w, color.Brown("-----> Info: "+message))
}
