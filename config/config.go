package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	GithubApp GithubConfig `yaml:"github"`
	AWS       AWSConfig    `yaml:"aws,omitempty"`
}

type GithubConfig struct {
	GithubAppIdentifier int    `yaml:"github-app-identifier"`
	GithubPrivateKey    string `yaml:"github-private-key"`
	GithubWebhookSecret string `yaml:"github-webhook-secret"`
}

// AWSConfig - only used for local development when access to aws api needed - not deployed
type AWSConfig struct {
	AWS_ACCESS_KEY_ID     string `yaml:"AWS_ACCESS_KEY_ID,omitempty"`
	AWS_SECRET_ACCESS_KEY string `yaml:"AWS_SECRET_ACCESS_KEY,omitempty"`
	AWS_SESSION_TOKEN     string `yaml:"AWS_SESSION_TOKEN,omitempty"`
	AWS_SECURITY_TOKEN    string `yaml:"AWS_SECURITY_TOKEN,omitempty"`
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
