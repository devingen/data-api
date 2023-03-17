# Start from a Debian image with Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.15

ADD . /go/src/github.com/devingen/data-api

WORKDIR /go/src/github.com/devingen/data-api

# Fetch dependencies and build the command inside the container.
RUN go install github.com/devingen/data-api/server

# Run the dock command by default when the container starts.
ENTRYPOINT server

# Service listens on port 1707
EXPOSE 80
