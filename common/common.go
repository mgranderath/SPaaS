package common

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/echo"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes generates a random string of length n
func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// HomeDir returns the home directory
func HomeDir() string {
	return os.Getenv("HOME")
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// CheckPasswordHash checks if a password hash is the password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Exists checks whether a file/directory exists
func Exists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

// EncodeJSONAndFlush encodes a struct to json and sends it
func EncodeJSONAndFlush(c echo.Context, response interface{}) error {
	if err := json.NewEncoder(c.Response()).Encode(response); err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}

// SpaasName returns name with a prefix such as to more easily define spaas only containers
func SpaasName(name string) string {
	return "spaas-" + name
}

// Copy copies a file
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
