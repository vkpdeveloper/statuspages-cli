package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	APIKey         string `api_key:"json"`
	ConfigFilePath string `config_file_path:"json"`
}

func InitAppConfig() (*AppConfig, error) {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, errors.New("failed to get the users home dir")
	}

	configFilePath := fmt.Sprintf("%s/.config/statuspages-cli.json", userHomeDir)

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("json")

	return &AppConfig{
		ConfigFilePath: configFilePath,
	}, nil
}

func (config *AppConfig) WriteConfig(apiKey string) bool {
	viper.Set("api_key", apiKey)

	viper.WriteConfig()

	return true
}

func (config *AppConfig) ReadConfig() bool {
	viper.ReadInConfig()
	apiKey := viper.GetString("api_key")

	if apiKey == "" {
		return false
	}

	config.APIKey = apiKey
	return true
}
