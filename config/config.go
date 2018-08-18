package config

import (
	"github.com/spf13/viper"
)

// ReadConfig reads a config file `filename` from `filePath` and uses the `defaults`
func ReadConfig(filePath string, filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(filePath)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
