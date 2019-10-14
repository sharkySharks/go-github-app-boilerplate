package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type AWSConfig struct {
	S3Bucket string `yaml:"s3Bucket"`
	S3Key    string `yaml:"s3Key"`
}
type Config struct {
	GithubApp GithubConfig `yaml:"github"`
}
type GithubConfig struct {
	GithubAppIdentifier int    `yaml:"github-app-identifier"`
	GithubPrivateKey    string `yaml:"github-private-key"`
	GithubWebhookSecret string `yaml:"github-webhook-secret"`
}

func retrieveSecretsFromS3() error {
	// check if a secrets.yaml file already exists and, if so, use this file and return from fn
	if _, err := os.Stat("secrets.yaml"); !os.IsNotExist(err) {
		log.Info("secrets.yaml file found, not connecting to s3...")
		return nil
	}
	log.Info("Connecting to s3...")
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	// get AWS S3 bucket name and key information from .aws/config.yaml file
	var awsConfig AWSConfig
	b, err := ioutil.ReadFile(".aws/config.yaml")
	if err != nil {
		return err
	}
	if err := yaml.UnmarshalStrict(b, &awsConfig); err != nil {
		return err
	}

	// the s3 bucket should be namespaced by environment, with the environment appended to the end of the bucket name
	env := os.Getenv("env")
	awsConfig.S3Bucket = awsConfig.S3Bucket + "-" + env

	// get S3 bucket/key data from S3
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsConfig.S3Bucket),
		Key:    aws.String(awsConfig.S3Key),
	})
	if err != nil {
		return err
	}
	secrets, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}

	// write S3 data to a secrets.yaml file
	writeError := ioutil.WriteFile("secrets.yaml", secrets, 0644)
	if writeError != nil {
		return writeError
	}
	log.Info("secrets.yaml file saved successfully!")
	return nil
}

func readConfig() (*Config, error) {
	var c Config

	// if `env=prod` or `env=stage`, then pull secrets from configured s3 bucket
	// configure AWS s3 bucket values in .aws/config.yaml file
	if env := os.Getenv("env"); env == "prod" || env == "stage" {
		err := retrieveSecretsFromS3()
		if err != nil {
			log.Fatal("Error retrieving secrets from S3: ", err)
		}
	}

	secretsPath := "secrets.yaml"
	bytes, err := ioutil.ReadFile(secretsPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
