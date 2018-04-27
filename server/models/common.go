package models

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"

	color "github.com/logrusorgru/aurora"
)

func piName(name string) string {
	return "pi-" + name
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FileExists returns if file exists
func FileExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString returns a random string of length n. Used to generate
// the secret
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetHomeFolder gets filepath to home folder of user
func GetHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// PrintErr prints an error
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

// PrintSuccess prints a success message
func PrintSuccess(w *os.File, message string) {
	fmt.Fprintln(w, color.Green("-----> Success: "+message))
}

// PrintNormal prints a normal message
func PrintNormal(w *os.File, message string) {
	fmt.Fprintln(w, ("-----> Task: " + message))
}

// PrintInfo prints a info message
func PrintInfo(w *os.File, message string) {
	fmt.Fprintln(w, color.Brown("-----> Info: "+message))
}
