SHELL := /bin/bash
UNAME := $(shell uname)

.PHONY: init_setup build-all build-app clean deploy gofmt run-local start

# only meant to run once per project setup, when first initializing the application from the boilerplate app code
# this script will:
# - checkout a new git branch based on the project name
# - customize the application files with your project name
# - change go.mod module to your repo location
# - add your repo as a new git remote to push the code to
init_setup: clean
	# binary is created based on mac or linux OS, change the GOOS and GOARCH env vars to match your system if you run into compile issues
	@if [ $(UNAME) == Darwin ]; then \
		export GOOS=darwin; \
	else \
		export GOOS=linux; \
	fi; \
	go mod tidy && \
	go mod download && \
	cd setup && \
	env GO111MODULE=on CGO_ENABLED=0 GOOS=$$GOOS GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/init_setup.go . && \
	cd ../ && ./bin/init_setup.go

# if you add more serverless functions then add another build step target and add it to the 'build-all' command
build-all: build-app

build-app: clean
	go mod tidy && \
	go mod download && \
	cd app && \
	env GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/app.go .

run-local:
	export CONFIG_FILE="./secrets.local.yaml" && npm start

start: build-app run-local

clean:
	rm -rf ./bin

deploy-stage: clean build-all
	export CONFIG_FILE="./secrets.stage.yaml" && sls deploy --verbose --stage stage

deploy-prod: clean build-all
	export CONFIG_FILE="./secrets.prod.yaml" && sls deploy --verbose --stage prod

gofmt:
	@echo "Formatting files..."
	gofmt -s -l -w config
	gofmt -s -l -w app
	gofmt -s -l -w setup

