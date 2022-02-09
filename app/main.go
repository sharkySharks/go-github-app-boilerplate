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
)

/*
	Add top level GitHub request payload keys to this struct based on the events the app is subscribing to and pull in
	the type from the go-github/github library: https://github.com/google/go-github
*/
type GitHubEvent struct {
	Action           string              `json:"action"`
	Installation     AppInstallation     `json:"installation,omitempty"`
	Issue            github.Issue        `json:"issue,omitempty"`
	IssueComment     github.IssueComment `json:"comment,omitempty"`
	Repo             github.Repository   `json:"repository,omitempty"`
	Sender           github.User         `json:"sender,omitempty"`
	XGithubEvent     string              `json:"-"` // this value comes from the header x-github-event
	XGithubRequestId string              `json:"-"` // this value comes from the header x-github-delivery
}

type AppInstallation struct {
	Id int64 `json:"id"`
}

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var (
	ghClient *github.Client
	conf     *config.Config
	event    *GitHubEvent
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

	/*
			if you need access to the aws api to take action based on an incoming event, then you will need credentials stored locally
		 	these values will only be used locally, they will not be deployed
		 	for deployment you will need to configure serverless.yaml to include the iam role statements with needed permissions
	*/
	//if os.Getenv("IS_OFFLINE") == "true" {
	//	log.Info("Setting AWS env vars in offline mode - local development only")
	//
	//	_ = os.Setenv("AWS_ACCESS_KEY_ID", conf.AWS.AWS_ACCESS_KEY_ID)
	//	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", conf.AWS.AWS_SECRET_ACCESS_KEY)
	//	_ = os.Setenv("AWS_SECURITY_TOKEN", conf.AWS.AWS_SECURITY_TOKEN)
	//	_ = os.Setenv("AWS_SESSION_TOKEN", conf.AWS.AWS_SESSION_TOKEN)
	//
	//	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
	//		log.Fatal("Need to set AWS credentials in secrets.local.yaml config file for local development")
	//		os.Exit(1)
	//	}
	//}
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
	log.Info(fmt.Sprintf("[%s] Payload validated and authenticated.", event.XGithubRequestId))
	res := app(request)
	return res, nil
}

func main() {
	lambda.Start(Handler)
}

// application
func app(request Request) Response {
	/*
		this section is where webhook events will be received after passing through middleware validation
		below is an example of handling a comment on a pull request including checking if a user has write permissions on the repo
		to work, the GitHub app would need to be configured in GitHub to subscribe to Issue Comment creation events
	*/
	if request.HTTPMethod == "POST" {
		switch event.XGithubEvent {
		case "issue_comment":
			if event.Action == "created" {
				// example code: check if the user has write permissions to the repo before allowing them to run tests
				hasPermission, err := hasWritePermission(event)
				if err != nil {
					err := fmt.Sprintf("[GitHub Request Id %s] error returned from GetPermissionLevel: %v", event.XGithubRequestId, err)
					log.Error(err)
					return Response{StatusCode: 424, Body: err}
				}
				if !hasPermission {
					msg := fmt.Sprintf("[GitHub Request Id %s] user unauthorized for running tests on repo. User: %s PR: %d",
						event.XGithubRequestId,
						*event.Sender.Login,
						*event.Issue.Number,
					)
					log.Info(msg)
					return Response{StatusCode: 200, Body: msg}
				}
				msg := fmt.Sprintf("[GitHub Request Id %s] User validated for repo with admin or write access: User %s, Repo: %s/%s",
					event.XGithubRequestId,
					*event.Sender.Login,
					*event.Repo.Owner.Login,
					*event.Repo.Name,
				)
				log.Info(msg)

				comment := strings.TrimSpace(*event.IssueComment.Body)
				switch comment {
				case "run all tests":
					log.Info(fmt.Sprintf("[GitHub Request Id %s] Received comment: %v", event.XGithubRequestId, comment))
					// execute some code here based on receiving a comment on a pull request
					// ie, respondToComment(comment, event)
					successMsg := fmt.Sprintf("Successfully received comment: %s", comment)
					_, _, err := ghClient.Issues.CreateComment(context.Background(), *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number,
						&github.IssueComment{
							Body: &successMsg,
						})
					if err != nil {
						log.Error(fmt.Sprintf("[GitHub Request Id %s] Error leaving comment on pull request %d: %v",
							event.XGithubRequestId,
							*event.Issue.Number,
							err,
						))
					}
					return Response{StatusCode: 200, Body: fmt.Sprintf("Received comment: %v", comment)}
				default:
					str := fmt.Sprintf("[GitHub Request Id %s] Received an unhandled comment: %s", event.XGithubRequestId, comment)
					log.Debug(str)
					return Response{StatusCode: 200, Body: str}
				}
			}
			fallthrough
		case "installation":
			if event.Action == "created" {
				log.Info(fmt.Sprintf("[GitHub Request Id %s] Received installation request", event.XGithubRequestId))
				return Response{StatusCode: 200, Body: "Received installation request"}
			}
			fallthrough
		default:
			e := fmt.Errorf("[GitHub Request Id %s] cannot find handler for event type: %s and/or action type: %s",
				event.XGithubRequestId,
				event.XGithubEvent,
				event.Action,
			)
			log.Error(e)
			return Response{Body: e.Error(), StatusCode: 200}
		}
	}

	e := fmt.Errorf("[GitHub Request Id %s] method not allowed: %s", event.XGithubRequestId, request.HTTPMethod)
	log.Error(e)
	return Response{Body: e.Error(), StatusCode: 405}
}

//hasPermission returns a boolean based on whether a user has the 'admin' or 'write' (true) or other permission level (false)
func hasWritePermission(event *GitHubEvent) (bool, error) {
	permissionLevel, _, err := ghClient.Repositories.GetPermissionLevel(context.Background(), *event.Repo.Owner.Login, *event.Repo.Name, *event.Sender.Login)
	if err != nil {
		return false, err
	}
	if *permissionLevel.Permission == "admin" || *permissionLevel.Permission == "write" {
		return true, nil
	}
	return false, nil
}
