package main

import (
	"encoding/json"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

var ghClient *github.Client

// validate payload webhook signature
func validatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := github.ValidatePayload(r, []byte(conf.GithubApp.GithubWebhookSecret))
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		pd, err := getWebhookAPIRequest(p)
		if err != nil {
			log.Error("Error unmarshalling json payload: ", err)
			return
		}
		payload = pd
		next.ServeHTTP(w, r)
	})
}

// helper function to convert json payload to struct
func getWebhookAPIRequest(body []byte) (*WebhookAPIRequest, error) {
	var wh = new(WebhookAPIRequest)
	err := json.Unmarshal(body, &wh)
	if err != nil {
		return nil, err
	}
	return wh, nil
}

// authenticate as Github App
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := initClient(payload.Installation.Id)
		if err != nil {
			log.Error("Error initializing client: ", err)
			return
		}
		ghClient = c

		next.ServeHTTP(w, r)
	})
}

func initClient(installationId int64) (*github.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(
		tr,
		conf.GithubApp.GithubAppIdentifier,
		installationId,
		[]byte(conf.GithubApp.GithubPrivateKey),
	)
	if err != nil {
		return nil, err
	}

	c := github.NewClient(&http.Client{Transport: itr})

	return c, nil
}
