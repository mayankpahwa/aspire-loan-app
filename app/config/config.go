package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the app config
type Config struct {
	Port          int `envconfig:"PORT" default:"8080"`
	DB            DBConfig
	MigrationPath string `envconfig:"MIGRATION_PATH"`
}

type DBConfig struct {
	DSN               string        `envconfig:"DB_DSN" required:"true"`
	MaxOpenConnection int           `envconfig:"MAX_OPEN_CONNECTION" default:"10"`
	MaxIdleConnection int           `envconfig:"MAX_IDLE_CONNECTION" default:"2"`
	MaxConnLifetime   time.Duration `envconfig:"MAX_CONNECTION_LIFETIME" default:"1m"`
}

// LoadConfig loads the configs from environment variable
func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
