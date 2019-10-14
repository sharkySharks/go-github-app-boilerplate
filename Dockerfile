FROM golang:1.12-alpine

ENV GOBIN $GOPATH/bin
RUN apk update && \
    apk add --virtual build-dependencies \
    git && \
    git config --global url."https://".insteadOf git://


WORKDIR $GOPATH/src/github.com/sharkySharks/go-github-app-boilerplate/

COPY . .

RUN cd main && \
    go get -d -v && \
    go install -v && \
    apk del build-dependencies && \
    rm -rf /var/cache/apk/*

EXPOSE 8080

ARG env=prod
ENV env=${env}

CMD ["/go/bin/main"]

