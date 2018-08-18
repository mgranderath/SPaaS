package config

import (
	"log"

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

// New creates a new store
func New(FilePath string, FileName string) {
	Cfg = &Store{
		FilePath: FilePath,
		FileName: FileName,
	}
	config, err := ReadConfig(FilePath, FileName, map[string]interface{}{
		"secret": RandStringBytes(15),
	})
	if err != nil {
		log.Panic("Error reading configuration file at " + FilePath + "/" + FileName)
	}
	Cfg.Config = config
}

// Save saves the configuration file
func Save() {
	if err := Cfg.Config.WriteConfig(); err != nil {
		log.Panic("Error writing configuration file at " + Cfg.FilePath + "/" + Cfg.FileName)
	}
}
