package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// middleware: validate payload webhook signature
func validatePayload(request Request) error {
	// depending on the proxy, headers can be lowercase or uppercase
	var (
		xHubSignature    string
		xGithubEvent     string
		xGitHubRequestId string
	)
	log.Info("Headers: ", request.Headers)

	if request.Headers["X-Hub-Signature"] == "" {
		xHubSignature = request.Headers["x-hub-signature"]
	} else {
		xHubSignature = request.Headers["X-Hub-Signature"]
	}
	p, err := github.ValidatePayloadFromBody(
		request.Headers["content-type"],
		strings.NewReader(request.Body),
		xHubSignature,
		[]byte(conf.GithubApp.GithubWebhookSecret))
	if err != nil {
		return err
	}

	pd, err := unmarshallGitHubRequest(p)
	if err != nil {
		log.Error("Error unmarshalling json payload: ", err)
		return err
	}
	event = pd
	// add header values to the event instance
	if request.Headers["X-GitHub-Event"] == "" {
		xGithubEvent = request.Headers["x-github-event"]
		xGitHubRequestId = request.Headers["x-github-delivery"]
	} else {
		xGithubEvent = request.Headers["X-GitHub-Event"]
		xGitHubRequestId = request.Headers["X-GitHub-Delivery"]
	}
	event.XGithubEvent = xGithubEvent
	event.XGithubRequestId = xGitHubRequestId

	log.Info(fmt.Sprintf("[GitHub Request Id %s] Event loaded.", event.XGithubRequestId))
	return nil
}

// middleware: helper function to convert json payload to struct
func unmarshallGitHubRequest(body []byte) (*GitHubEvent, error) {
	var event = new(GitHubEvent)
	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// middleware: authenticate as Github App
func authenticate() error {
	log.Info(fmt.Sprintf("authenticating request for event: %v", event))
	c, err := initGitHubClient(
		event.Installation.Id,
		int64(conf.GithubApp.GithubAppIdentifier),
		[]byte(conf.GithubApp.GithubPrivateKey),
	)
	if err != nil {
		log.Error("Error initializing client: ", err)
		return err
	}
	ghClient = c
	return nil
}

// initGitHubClient initializes a github client to respond back to github after events are received
func initGitHubClient(installationId int64, githubAppID int64, githubPrivateKey []byte) (*github.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(
		tr,
		githubAppID,
		int64(installationId),
		githubPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	c := github.NewClient(&http.Client{Transport: itr})

	return c, nil
}
