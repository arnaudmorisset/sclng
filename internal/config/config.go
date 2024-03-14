package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type GithubConfig struct {
	Token string `envconfig:"GITHUB_TOKEN"`
}

type Config struct {
	Port   int `envconfig:"PORT" default:"8080"`
	Github GithubConfig
}

func NewConfig() (Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fail to process the environment variables: %s", err.Error())
	}

	return cfg, nil
}
