package configuration

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API APIConfig
}

func LoadConfig() (*Config, error) {
	var config Config
	noPrefix := ""
	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type APIConfig struct {
	Port string `envconfig:"API_PORT" default:"3000"`
}
