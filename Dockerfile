FROM golang:1.17-alpine as builder

ENV GOBIN $GOPATH/bin
RUN apk update && \
    apk add --virtual build-dependencies \
    git && \
    git config --global url."https://".insteadOf git://


WORKDIR $GOPATH/src/github.com/sharkySharks/go-github-app-boilerplate/

COPY . .

RUN cd main && \
    go mod tidy && \
    go install && \ 
    cp secrets.yaml /secrets.yaml

EXPOSE 8080

ARG env=prod
ENV env=${env}

FROM alpine:3.12

# Copy our static executable.
COPY --from=builder /go/bin/main /bin/main
COPY --from=builder /secrets.yaml /secrets.yaml

CMD ["/bin/main"]