# go-github-app-boilerplate - [Container version - not being maintained in favor of serverless]

![go-github](go-github.png)

Boilerplate for creating a GitHub App in Golang. This is the `container` version of the application setup. If you would like to run a github application using [serverless and AWS lambdas](https://www.serverless.com/) then checkout the `serverless` branch.

[GitHub apps](https://docs.github.com/en/free-pro-team@latest/developers/apps/getting-started-with-apps) are APIs that are configured with specific GitHub API credentials so that the API can receive and authenticate payloads from GitHub events.

## Getting started

This repository is written in [golang](https://golang.org/), aka `go`. If you have not installed go before and you wish to make contributions to this repository, then follow these [installation instructions](https://golang.org/doc/install) before proceeding.
This code uses golang version >=1.16 and [go modules](https://go.dev/blog/using-go-modules).

You can also run this application using Docker, and therefore do not need to install golang.

### GitHub Application Setup and Secrets Configuration

To create your own GitHub Application to use with this repository:

1. Create a GitHub application on GitHub, following these [instructions](https://developer.github.com/apps/building-github-apps/creating-a-github-app/).
This [link](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/) may also be helpful when setting up a new GitHub application.

2. Clone this repository and copy `secrets.example.yaml` to `secrets.stage.yaml` and `secrets.prod.yaml` in the root of the repository. Fill in the values in `secrets.*.yaml` with respective data for your GitHub application and for that environment. See table below for `secrets.*.yaml` key-value mappings.

#### secrets.*.yaml

| Key                      | Description                                                              | Default/Type                     |
|--------------------------|--------------------------------------------------------------------------| ---------------------------------|
| `github-app-identifier`  | The ID of the application under the _About_ section found under `Settings / Developer settings / GitHub Apps / your-app` | `None`; type: int |
| `github-webhook-secret`  | The webhook secret created when the application was created. This can be reset under the settings for the application. | `None`; type string |
| `github-private-key`     | The private key generated during application creation, this can also be reset under the settings for the application. | `None`; plaintext RSA key |

### Running This Repo Code Locally

#### Run Locally With Docker:

1. Make sure you have [Docker installed](https://docs.docker.com/v17.12/install/).

2. Check that you have a `secrets.stage.yaml` file with your GitHub application configuration values in the root of the repository. This secrets file will be used for local development.

3. From the root of the repository, run `docker build --build-arg env=stage -t my-app:latest .` to build the Docker image.

4. After the image has finished building, run `docker run -e env=stage -p 8080:8080 my-app:latest` and you should see a server listening message in your terminal output. This server is set up to receive webhook requests for the application configured in `secrets.stage.yaml`.

5. Get your smee.io link that you [setup earlier](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/#step-1-start-a-new-smee-channel), and run that in another terminal window. Example: `smee --url https://smee.io/qrfeVRbFbffd6vD --path / --port 8080`

6. You are now able to test that you are receiving webhook requests from GitHub. Test this out by sending an event from GitHub and see the output of the request in both terminal windows. The default app is setup to handle comments made on a pull request.


### Pulling Your Secrets From An AWS S3 Bucket

This application supports pulling secrets from an s3 bucket, for both local development as well as for stage/prod deployments.

* First, store your `secrets.yaml` file in an encrypted s3 bucket with `secrets.yaml` as the key and `-stage` or `-prod` at the end of your bucket name, depending on the environment. s3 buckets are globally scoped, so the names must be unique across all AWS accounts.

* Second, add your AWS S3 bucket name and key to the `.aws/config.yaml` file, leaving off the `-stage` or `-prod` from the bucket name, this will be added based on the `env` environment variable set (see Third step).

* Third, set the `env` environment variable to the correct environment, either `stage` for staging or `prod` for production in order to pull from a specific s3 bucket. The environment (`stage` or `prod`) should be appended to the end of the bucket name, ie if the bucket name is `my-app-secrets` then the expectation is that there is a bucket named `my-app-secrets-stage` and/or `my-app-secrets-prod` in those respective AWS environments.

#### Local Setup

1. Make sure to configure your AWS credentials, [see docs](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).


##### Run With Docker

* From the root of the repository, build the docker image and pass in the following environment variables to the run command:
```
docker run
-e AWS_ACCESS_KEY_ID=<your-aws-access-key>
-e AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key>
-e AWS_REGION=<your-aws-region>
-p 8080:8080 mts:latest
```

| Environment Variables           | Description                                                              | Default |
|---------------------------------|--------------------------------------------------------------------------|---------|
| `AWS_ACCESS_KEY_ID`             | AWS credential - AWS Access Key Id                                       | `None`  |
| `AWS_SECRET_ACCESS_KEY`         | AWS credential - AWS Secret Access Key                                   | `None`  |
| `AWS_REGION`                    | The default AWS region where the s3 bucket lives                         | `None`  |


#### Production/Multi Environment Setup

* The Dockerfile's default `env` is set to `stage`, but you can pass the environment variable as a [build-arg](https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg), ie `--build-arg env=prod`.
In order for your deployment service to be able to access your s3 bucket, you will need to add the `AmazonS3ReadOnlyAccess` policy to your service role, as well as add a policy to your s3 bucket that gives permission for your service role to access your s3 bucket.

* Once these are setup, just make sure that you have your S3 bucket information saved in `.aws/config.yaml`, following the steps at the beginning of this section, titled: `Pulling Your Secrets From An s3 Bucket`, or set the environment variables in your shell window.

## Future Features

"...a work is never truly completed [...] but abandoned..." Paul Valéry

Nice-to-have features:
- add testing
