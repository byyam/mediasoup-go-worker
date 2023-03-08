#!/bin/bash

set -x

go version

export GO111MODULE=on

CGO_ENABLED=0 go build -mod=mod -o bin/example-uds-controller example/uds/controller/main.go
CGO_ENABLED=0 go build -mod=mod -o bin/example-uds-controlled example/uds/controlled/main.go
