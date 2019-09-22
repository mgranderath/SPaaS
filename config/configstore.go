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

// Cfg points to the current config
var Cfg *Store

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
func New(FilePath string, FileName string) {
	Cfg = &Store{
		FilePath: FilePath,
		FileName: FileName,
	}
	_, _ = os.OpenFile(FilePath+"/"+FileName, os.O_RDONLY|os.O_CREATE, 0666)
	err := os.MkdirAll(filepath.Join(common.HomeDir(), ".spaas", "acme"), os.ModePerm)
	if err != nil {
		log.Fatalln("Could not create acme folder")
	}
	err = os.MkdirAll(filepath.Join(common.HomeDir(), ".spaas", "applications"), os.ModePerm)
	if err != nil {
		log.Fatalln("Could not create applications folder")
	}
	config, err := ReadConfig(FilePath+"/"+FileName, defaultConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	config.AutomaticEnv()
	Cfg.Config = config
}

// Save saves the configuration file
func Save() error {
	if err := Cfg.Config.WriteConfigAs(Cfg.FilePath + "/" + Cfg.FileName); err != nil {
		return err
	}
	return nil
}
