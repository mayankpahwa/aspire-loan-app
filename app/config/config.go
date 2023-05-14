package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the app config
type Config struct {
	Port int `envconfig:"PORT" default:"8080"`
}

// LoadConfig loads the configs from environment variable
func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
