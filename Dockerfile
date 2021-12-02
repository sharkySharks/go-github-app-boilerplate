FROM golang:1.17-alpine as builder

ENV GOBIN $GOPATH/bin
ARG env=stage
ENV env=${env}

RUN apk update && \
    apk add --virtual build-dependencies \
    git && \
    git config --global url."https://".insteadOf git://

WORKDIR $GOPATH/src/github.com/sharkySharks/go-github-app-boilerplate/

ADD ./.aws ./main ./secrets.${env}.yaml ./

RUN go mod tidy && \
    go install


FROM alpine:3.12

ARG env=stage
ENV env=${env}

# Copy our static executable.
COPY --from=builder /go/bin/go-github-app-boilerplate /bin/go-github-app-boilerplate
# Copy the config to the root.
COPY --from=builder /go/src/github.com/sharkySharks/go-github-app-boilerplate/secrets.${env}.yaml /secrets.yaml

CMD ["/bin/go-github-app-boilerplate"]

EXPOSE 8080
