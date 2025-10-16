package box

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HostName    string `envconfig:"GO_BASE_HOST_NAME"`
	ServiceName string `envconfig:"GO_BASE_SERVICE_NAME"`
	Debug       bool   `envconfig:"GO_BASE_DEBUG" default:"false"`

	PostgresConfig PostgresConfig
}

type PostgresConfig struct {
	ConnectionString string `envconfig:"GO_BASE_POSTGRES_CONNECTION_STRING" required:"true"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}

	if hostName, err := os.Hostname(); err == nil {
		cfg.HostName = hostName
	}

	if err := envconfig.Process("GO_BASE", cfg); err != nil {
		return nil, fmt.Errorf("error while parsing environment variables | %w", err)
	}

	return cfg, nil
}
