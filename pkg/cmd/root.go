package cmd

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	path     string
	fileType string
	name     string
}

func NewConfigYaml(path string, name string) *Config {
	return &Config{path: path, name: name, fileType: "yaml"}
}

// InitConfig reads in Config file and ENV variables if set.
func InitConfig(cfg *Config) {
	viper.AddConfigPath(cfg.path)
	viper.SetConfigType(cfg.fileType)
	viper.SetConfigName(cfg.name)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	log.Printf("Using config file: %s", viper.ConfigFileUsed())
}
