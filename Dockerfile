# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/andrewm3/pe-code-mgr-webhook-adapter

# Allow binary to perform HTTPS requests
RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

# Build the pe-code-mgr-webhook-adapter command inside the container.
RUN go install github.com/andrewm3/pe-code-mgr-webhook-adapter

# Run the pe-code-mgr-webhook-adapter command by default when the container starts.
ENTRYPOINT /go/bin/pe-code-mgr-webhook-adapter

# Document that the service listens on port 8080.
EXPOSE 8080
