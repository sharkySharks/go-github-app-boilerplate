.PHONY: build-all build-app clean deploy gofmt run-local start

# if you add more serverless functions then add another build step target and add it to the 'build-all' command
build-all: build-app

build-app: clean
	go mod tidy && \
	go mod download && \
	cd app && \
	export GO111MODULE=on && \
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/app.go .

run-local:
	export CONFIG_FILE="./secrets.stage.yaml" && npm start

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
	gofmt -s -l -w github
	gofmt -s -l -w app

