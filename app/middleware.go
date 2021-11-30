package main

import (
	"encoding/json"
	"github.com/google/go-github/v37/github"
	gc "github.com/sharkysharks/go-github-app-boilerplate/github"
	log "github.com/sirupsen/logrus"
	"strings"
)

// middleware: validate payload webhook signature
func validatePayload(request Request) error {
	// depending on the proxy, headers can be lowercase or uppercase
	var (
		xHubSignature string
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

	pd, err := getWebhookAPIRequest(p)
	if err != nil {
		log.Error("Error unmarshalling json payload: ", err)
		return err
	}
	payload = pd
	return nil
}

// middleware: helper function to convert json payload to struct
func getWebhookAPIRequest(body []byte) (*WebhookAPIRequest, error) {
	var wh = new(WebhookAPIRequest)
	err := json.Unmarshal(body, &wh)
	if err != nil {
		return nil, err
	}
	return wh, nil
}

// middleware: authenticate as Github App
func authenticate() error {
	log.Info("authenticating request for payload: ", payload)
	c, err := gc.InitClient(
		payload.Installation.Id,
		int64(conf.GithubApp.GithubAppIdentifier),
		[]byte(conf.GithubApp.GithubPrivateKey),
	)
	if err != nil {
		log.Error("Error initializing client: ", err)
		return err
	}
	client = c
	return nil
}
