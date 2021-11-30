# go-github-app-boilerplate

![go-github](go-github.png)

Boilerplate for creating a GitHub App in Golang. This is the serverless/lambda version. For the container version of the application, checkout the `container` branch.

[GitHub apps](https://docs.github.com/en/free-pro-team@latest/developers/apps/getting-started-with-apps) are APIs that are configured with specific GitHub API credentials so that the API can receive and authenticate payloads from GitHub events.
For a list of github event types that can be subscribed to and acted upon, see github docs: https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types

This application uses Golang version >=1.16 and [go modules](https://go.dev/blog/using-go-modules).

## GitHub Application Setup and Secrets Configuration

To create your own GitHub Application to use with this repository:

1. Create a GitHub application on GitHub, following these [instructions](https://developer.github.com/apps/building-github-apps/creating-a-github-app/).
This [link](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/) may also be helpful when setting up a new GitHub application.

2. Clone this repository and copy `secrets.example.yaml` to `secrets.stage.yaml` and `secrets.prod.yaml` in the root of the repository. Fill in the values in `secrets.*.yaml` with respective data for your GitHub application. See table below for `secrets.yaml` key-value mappings.

### secrets.yaml

| Key                      | Description                                                              | Default/Type                     |
|--------------------------|--------------------------------------------------------------------------| ---------------------------------|
| `github.github-app-identifier`  | The ID of the application under the _About_ section found under `Settings / Developer settings / GitHub Apps / your-app` | `None`; type: int |
| `github.github-webhook-secret`  | The webhook secret created when the application was created. This can be reset under the settings for the application. | `None`; type string |
| `github.github-private-key`     | The private key generated during application creation, this can also be reset under the settings for the application. | `None`; plaintext RSA key |

## Local Development
This application is set up as a lambda function and uses the [serverless](https://www.serverless.com/) framework to develop locally and deploy.

Make sure to [install serverless](https://www.serverless.com/framework/docs/getting-started/) on your computer.

It is written in [Golang version 1.16](https://golang.org/doc/install), so make sure you also have this installed.
 
The application expects two secrets files for deployment: `secrets.stage.yaml` and `secrets.prod.yaml`.
Go back to the previous section if you have not configured this secrets file yet.
For local development `secrets.stage.yaml` is used. 

After installing the above you should be able to run the following commands to start the application:

```bash
npm install
make start
```

These commands will build the golang binary and run the lambda through the [serverless-offline plugin](https://github.com/dherault/serverless-offline).

For extra debugging output, add `SLS_DEBUG=*` in front of the `npm > start` command in `package.json` like so: `"start": "SLS_DEBUG=* sls offline start --printOutput"`.

To receive events from the GitHub staging application, visit the application you set up in github.com and set the webhook url to a smee proxy url, which you can create at [smee.io](smee.io).

Run the smee proxy locally with the following command: 
```
smee --url https://smee.io/c926vE5gmuwgsGY --path /webhooks --port 3000
```
This will forward all traffic received by the smee url to `localhost:3000/`, which is the default location of the locally running serverless application.

Test this out by leaving a comment on a pull request in a repo that you installed the github application on in github.com. You should see the event be captured in the `Advanced` tab in the github application settings page in github.com, as well as the event received by the smee proxy and the application.


## Deployment - Stage/Prod
After working locally, you can start testing in the staging environment and then eventually deploy to production.

The serverless framework deploys lambda functions into AWS, therefore you should make sure that your [AWS credentials are configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) locally before running the following commands:

```bash
# deploy to staging env:
make deploy-stage

# deploy to prod env:
make deploy-prod
```

Once the lambda is deployed, update the webhook url in the relevant github application configuration, noted above, in github.com. This will then start sending the requests to your lamdba function.

Check the logs in the lambda function as well as on the github application event stream (under `Advanced` side panel).

Note: you may see some requests receive a *time out* error in GitHub, but you will see that the request actually did complete, the process just took longer than one service was expecting.

### Remove Resources
If you want to remove something that you created, then have your AWS credentials set and run `sls remove`. This will remove all the resources that the serverless framework created.

## Future Features

"...a work is never truly completed [...] but abandoned..." Paul Valéry

Nice-to-have features:
- add testing
