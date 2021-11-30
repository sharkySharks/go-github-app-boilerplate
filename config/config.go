package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	GithubApp GithubConfig `yaml:"github"`
}

type GithubConfig struct {
	GithubAppIdentifier int    `yaml:"github-app-identifier"`
	GithubPrivateKey    string `yaml:"github-private-key"`
	GithubWebhookSecret string `yaml:"github-webhook-secret"`
}

func ReadConfig(path string) (*Config, error) {
	var c Config

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
