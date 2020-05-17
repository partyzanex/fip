package env

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	LogLevel string `envconfig:"log_level" default:"debug"`
	Address  string `envconfig:"address" default:"localhost:9090"`
	Source   string `envconfig:"source"`
	Cache    string `envconfig:"cache" default:"/tmp/fip"`
}

func Read(prefix string) (*Config, error) {
	config := Config{}

	err := envconfig.Process(prefix, &config)
	if err != nil {
		return nil, errors.Wrap(err, "read config from environment failed")
	}

	return &config, nil
}
