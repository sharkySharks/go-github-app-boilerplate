package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v37/github"
	"github.com/sharkysharks/go-github-app-boilerplate/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

/*
	Add top level GitHub request payload keys to this struct based on the events the app is subscribing to and pull in
	the type from the go-github/github library
*/
type WebhookAPIRequest struct {
	Action       string              `json:"action,omitempty"`
	Installation AppInstallation     `json:"installation,omitempty"`
	Issue        github.Issue        `json:"issue,omitempty"`
	IssueComment github.IssueComment `json:"comment,omitempty"`
	Repo         github.Repository   `json:"repository,omitempty"`
	Sender       github.User         `json:"sender,omitempty"`
}

type AppInstallation struct {
	Id int64 `json:"id"`
}

type Event struct {
	Type   string
	Repo   string
	Action string
	Sender string
	ID     string
	Num    int64
}

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var (
	client  *github.Client
	conf    *config.Config
	mutex   sync.Mutex
	payload *WebhookAPIRequest
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Fatal("Could not locate configuration file. Set CONFIG_PATH environment variable to file location.")
	}
	c, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	conf = c
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	err := validatePayload(request)
	if err != nil {
		log.Errorf("Error validating payload: %v", err)
		return Response{StatusCode: 401}, err
	}
	err = authenticate()
	if err != nil {
		log.Errorf("Error authenticating request: %v", err)
		return Response{StatusCode: 401}, err
	}
	log.Info("Payload validated and authenticated.")
	res := app(request)
	return res, nil
}

func main() {
	lambda.Start(Handler)
}

// application
func app(request Request) Response {
	/*
		xGithubEvent : github event that is being received. See github docs for more info on event types: https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types
		xGithubDelivery : match the xGithubDelivery to the webhook event logs under the github app settings in the github web console
						  this is helpful for debugging webhook events. depending on the proxy being used, headers can be lowercase or uppercase
		event : add specific payload data to the event struct to pass to your next function
	*/
	var (
		xGithubEvent    string
		xGithubDelivery string
	)
	if request.Headers["X-GitHub-Event"] == "" {
		xGithubEvent = request.Headers["x-github-event"]
		xGithubDelivery = request.Headers["x-github-delivery"]
	} else {
		xGithubEvent = request.Headers["X-GitHub-Event"]
		xGithubDelivery = request.Headers["X-GitHub-Delivery"]
	}

	var event = Event{
		xGithubEvent,
		*payload.Repo.Name,
		payload.Action,
		*payload.Sender.Login,
		xGithubDelivery,
		*payload.Issue.ID,
	}

	log.Info("app:xGithubEvent: ", xGithubEvent)
	log.Info("app:event: ", event)

	/*
		this section is where webhook events will be received after passing through middleware validation
		below is an example of handling a comment on a pull request
		to work, the GitHub app would need to be configured in GitHub to subscribe to Issue Comment creation events
	*/
	if request.HTTPMethod == "POST" {
		switch xGithubEvent {
		case "issue_comment":
			if payload.Action == "created" {
				comment := strings.TrimSpace(*payload.IssueComment.Body)
				switch comment {
				case "run all tests":
					log.Info("Received event to run all tests")
					// execute some code here based on receiving a comment on a pull request
					// ie, respondToComment(comment, event)
					return Response{StatusCode: 200, Body: fmt.Sprintf("Received comment: %v", comment)}
				default:
					str := fmt.Sprintf("Received an unhandled comment: %v", comment)
					log.Info(str)
					return Response{StatusCode: 404, Body: str}
				}
			}
			fallthrough
		case "installation":
			if payload.Action == "created" {
				log.Info("Received installation request")
				return Response{StatusCode: 200, Body: "Received installation request"}
			}
			fallthrough
		default:
			e := fmt.Errorf("cannot find handler for event type: %v and/or action type: %v", xGithubEvent, payload.Action)
			log.Error(e)
			return Response{Body: e.Error(), StatusCode: 404}
		}
	} else {
		e := fmt.Errorf("method not allowed: %v", request.HTTPMethod)
		log.Error(e)
		return Response{Body: e.Error(), StatusCode: 405}
	}
}
