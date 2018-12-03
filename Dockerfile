# Use a build stage to minimise image size
FROM golang:1.11-alpine as build

# Copy the local package files to the container's workspace.
WORKDIR /go/src/github.com/andrewm3/pe-code-mgr-webhook-adapter
COPY . /go/src/github.com/andrewm3/pe-code-mgr-webhook-adapter

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

# The binary stage
FROM alpine:latest

# Allow binary to perform HTTPS requests
RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

# Copy the binary
WORKDIR /root/
COPY --from=build /go/src/github.com/andrewm3/pe-code-mgr-webhook-adapter/app .

# Specify default Code Manager endpoint to be localhost
ENV CODE_MGR_URL https://localhost:8170/code-manager/v1/webhook

# Run the pe-code-mgr-webhook-adapter command by default when the container starts.
ENTRYPOINT ./app --redirect "${CODE_MGR_URL}"

# Document that the service listens on port 8080.
EXPOSE 8080
