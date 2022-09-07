#!/bin/bash

set -x

go version

export GO111MODULE=on

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.githash=$(git rev-parse HEAD)" -mod=mod -o bin/mediasoup-worker-linux cmd/mediasoup-worker/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.githash=$(git rev-parse HEAD)" -mod=mod -o bin/mediasoup-worker-darwin cmd/mediasoup-worker/main.go
