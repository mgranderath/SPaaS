package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mgranderath/SPaaS/common"
	"github.com/spf13/viper"
)

// Store contains the current configuration
type Store struct {
	Config   *viper.Viper
	FilePath string
	FileName string
}

var defaultConfig = map[string]interface{}{
	"secret":            common.RandStringBytes(15),
	"username":          "spaas",
	"password":          common.HashPassword("smallpaas"),
	"letsencrypt":       false,
	"letesencryptEmail": "example@example.com",
	"acmePath":          filepath.Join(common.HomeDir(), ".spaas-server", "acme"),
	"domain":            "example.com",
	"useDomain":         false,
}

// New creates a new store
func New() *Store {
	filePath := filepath.Join(common.HomeDir(), ".spaas")
	fileName := ".spaas.json"
	ConfigStore := &Store{
		FilePath: filePath,
		FileName: fileName,
	}
	_, _ = os.OpenFile(filePath+"/"+fileName, os.O_RDONLY|os.O_CREATE, 0666)
	err := os.MkdirAll(filepath.Join(common.HomeDir(), ".spaas", "acme"), os.ModePerm)
	if err != nil {
		log.Fatalln("Could not create acme folder")
	}
	err = os.MkdirAll(filepath.Join(common.HomeDir(), ".spaas", "applications"), os.ModePerm)
	if err != nil {
		log.Fatalln("Could not create applications folder")
	}
	config, err := ReadConfig(filePath+"/"+fileName, defaultConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	config.AutomaticEnv()
	ConfigStore.Config = config
	if err := ConfigStore.Save(); err != nil {
		fmt.Println(err.Error())
	}
	return ConfigStore
}

// Save saves the configuration file
func (store *Store) Save() error {
	if err := store.Config.WriteConfigAs(store.FilePath + "/" + store.FileName); err != nil {
		return err
	}
	return nil
}
