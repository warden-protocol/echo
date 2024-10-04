package main

import (
	"errors"
	"fmt"

	env "github.com/caarlos0/env/v10"
)

//nolint:lll
type Config struct {
	Port            string   `env:"PORT"             envDefault:"10010"                  mapstructure:"PORT"`
	Endpoints       []string `env:"ENDPOINTS"        envDefault:"http://localhost:26657" mapstructure:"ENDPOINTS"        envSeparator:","`
	Peers           []string `env:"PEERS"                                                mapstructure:"PEERS"            envSeparator:","`
	BehindThreshold int64    `env:"BEHIND_THRESHOLD" envDefault:"10"                     mapstructure:"BEHIND_THRESHOLD"`
}

var errConfig = errors.New("config error")

func configError(msg string) error {
	return fmt.Errorf("%w: %s", errConfig, msg)
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	var err error

	// setDefaults(*cfg)

	if err = env.Parse(&cfg); err != nil {
		return Config{}, configError(err.Error())
	}

	return cfg, nil
}
