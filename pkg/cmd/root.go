package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
