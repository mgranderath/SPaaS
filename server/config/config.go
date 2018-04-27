package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Configuration Stores the main configuration for the application
type Configuration struct {
	Nginx  bool   `json:"nginx"`
	Secret string `json:"secret"`
}

// ReadConfig will read the configuration json file to read the parameters
// which will be passed in the config file
func ReadConfig(fileName string) (Configuration, error) {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return Configuration{}, err
	}
	var config Configuration
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print("Invalid JSON")
		return Configuration{}, err
	}
	return config, nil
}

// WriteConfig will write a configuratiuon to a file
func WriteConfig(fileName string, config Configuration) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(fileName, configJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}
