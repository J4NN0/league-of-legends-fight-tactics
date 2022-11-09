package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RiotAPIKey string `envconfig:"RIOT_API_KEY"`
	LoLRegion  string `envconfig:"LOL_REGION" required:"true"`
}

func ReadConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
