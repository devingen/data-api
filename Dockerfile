# Start from a Debian image with Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.22

ADD . /go/src/github.com/devingen/data-api

WORKDIR /go/src/github.com/devingen/data-api

# Fetch dependencies and build the command inside the container.
RUN go install github.com/devingen/data-api/cmd/api

ENTRYPOINT api
EXPOSE 80
