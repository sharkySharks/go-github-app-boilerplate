# go-github-app-boilerplate
Boilerplate for creating a GitHub Application in Golang

## Getting started

This repository is written in [golang](https://golang.org/), aka `go`. If you have not installed go before and you wish to make contributions to this repository, then follow these [installation instructions](https://golang.org/doc/install) before proceeding.
[How to Write Go Code](https://golang.org/doc/code.html) is a great starter document to help set up your environment properly.

You can also run this application using Docker, and therefore do not need to install golang.

### Expected directory structure

Following the guidelines of [How to Write Go Code](https://golang.org/doc/code.html), your directory structure should mirror the following:

 ```
 .
 └── go 
    ├── bin
    ├── pkg
    └── src
      └── github.com
          └── sharkySharks
             └── go-github-app-boilerplate
                ├── .aws
                │   └── config.yaml
                ├── Dockerfile
                ├── README.md
                ├── main
                │   ├── config.go
                │   ├── main.go
                │   └── middleware.go
                ├── vendor/               <- has all the vendor packages/dependencies
                ├── Gopkg.lock            <- dependency lock file
                ├── Gopkg.toml            <- dependency management file
                ├── secrets.example.yaml
                └── secrets.yaml
```

### GitHub Application Setup and Secrets Configuration

To create your own GitHub Application to use with this repository:

1. Create a GitHub application on GitHub, following these [instructions](https://developer.github.com/apps/building-github-apps/creating-a-github-app/).
This [link](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/) may also be helpful when setting up a new GitHub application.

2. Clone this repository and copy `secrets.example.yaml` to `secrets.yaml` in the root of the repository. Fill in the values in `secrets.yaml` with respective data for your GitHub application. See table below for `secrets.yaml` key-value mappings.

#### secrets.yaml

| Key                      | Description                                                              | Default/Type                     |
|--------------------------|--------------------------------------------------------------------------| ---------------------------------|
| `github-app-identifier`  | The ID of the application under the _About_ section found under `Settings / Developer settings / GitHub Apps / your-app` | `None`; type: int |
| `github-webhook-secret`  | The webhook secret created when the application was created. This can be reset under the settings for the application. | `None`; type string |
| `github-private-key`     | The private key generated during application creation, this can also be reset under the settings for the application. | `None`; plaintext RSA key |

### Running This Repo Code Locally

This project uses [dep](https://golang.github.io/dep/) to manage dependencies. See docs to install.

#### Run Locally With Docker:

1. Make sure you have [Docker installed](https://docs.docker.com/v17.12/install/).

2. Check that you have a `secrets.yaml` file with your GitHub application configuration values in the root of the repository.

3. From the root of the repository, run `docker build -t my-app:latest .` to build the Docker image.

4. After the image has finished building, run `docker run -e env=dev -p 8080:8080 my-app:latest` and you should see a server listening message in your terminal output. This server is set up to receive webhook requests for the application configured in `secrets.yaml`.

5. Get your smee.io link that you [setup earlier](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/#step-1-start-a-new-smee-channel), and run that in another terminal window. Example: `smee --url https://smee.io/qrfeVRbFbffd6vD --path / --port 8080`

6. You are now able to test that you are receiving webhook requests from GitHub. Test this out by sending an event from GitHub and see the output of the request in both terminal windows.

#### Run On Local File System:

1. Make sure you have [Golang installed](https://golang.org/doc/install).

2. Check that you have a `secrets.yaml` file with your GitHub application configuration values in the root of the repository.

3. From inside the `main/` directory, run `go get && go install`, which will create the executable file `main` in the `$GOBIN` directory.

4. From the root of the repository, run `$GOBIN/main` and you should see a server listening message. This server is set up to receive webhook requests for the application configured in `secrets.yaml`.

5. Get your smee.io link that you [setup earlier](https://developer.github.com/apps/quickstart-guides/setting-up-your-development-environment/#step-1-start-a-new-smee-channel), and run that in another terminal window. Example: `smee --url https://smee.io/qrfeVRbFbffd6vD --path / --port 8080`

6. You are now able to test that you are receiving webhook requests from GitHub. Test this out by sending an event that the GitHub App is configured to listen for and see the output of the request in both terminal windows.


### Pulling Your Secrets From An AWS S3 Bucket

This application supports pulling secrets from an s3 bucket, for both local development as well as for prod deployments.

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

##### Run On Local File System

* If you have configured your AWS credentials in step 1, then you should have a `~/.aws/credentials` file that looks like the following, though you may need to add a default region:

```
[default]
aws_access_key_id = <your-access-key-id>
aws_secret_access_key = <your-access-key>
region = <your-region>
```

* The AWS CLI will look up the credential chain and find this file. You may also set these values in your shell window.

* To run the github app, simply run `go get && go install` from inside the `main/` directory,

* Then run `env=prod $GOBIN/main` or `env=stage $GOBIN/main` to trigger pulling secrets from your specified s3 bucket.

#### Production/Multi Environment Setup

* The Dockerfile's default `env` is set to `prod`, but you can also pass the environment variable as a [build-arg](https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg), ie `--build-arg env=stage`.
In order for your deployment service to be able to access your s3 bucket, you will need to add the `AmazonS3ReadOnlyAccess` policy to your service role, as well as add a policy to your s3 bucket that gives permission for your service role to access your s3 bucket.

* Once these are setup, just make sure that you have your S3 bucket information saved in `.aws/config.yaml`, following the steps at the beginning of this section, titled: `Pulling Your Secrets From An s3 Bucket`, or set the environment variables in your shell window.

