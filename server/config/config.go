package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Configuration Stores the main configuration for the application
type Configuration struct {
	Nginx bool `json:"nginx"`
}

var err error
var config Configuration

//ReadConfig will read the configuration json file to read the parameters
//which will be passed in the config file
func ReadConfig(fileName string) (Configuration, error) {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return Configuration{}, err
	}
	//log.Print(configFile)
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print("Invalid JSON")
		return Configuration{}, err
	}
	return config, nil
}
