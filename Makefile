GOOS=darwin
ifeq ($(os),linux)
	GOOS=linux
endif

MOD=mod
ifeq ($(mod),vendor)
	MOD=vendor
endif

GOARCH=arm64
ifeq ($(arch),amd64)
	GOARCH=amd64
endif

BUILD_DATE=`date '+%Y-%m-%d_%H:%M:%S'`
GIT_HASH=`git rev-parse HEAD`
GIT_BRANCH=`git branch --show-current`
GO_VERSION=`go version|awk '{print $$3}'|sed 's/[ ][ ]*/_/g'`
LDFLAGS=-ldflags "-X=main.buildstamp=${BUILD_DATE} -X=main.githash=${GIT_HASH} -X=main.gitbranch=${GIT_BRANCH} -X=main.goversion=${GO_VERSION}"

all: mediasoup-worker sfu-server

mediasoup-worker:
	@go version
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build $(LDFLAGS) -mod=${MOD} -o bin/mediasoup-worker-${GOOS} cmd/mediasoup-worker/main.go

sfu-server:
	@go version
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build $(LDFLAGS) -mod=${MOD} -o bin/sfu-server-${GOOS} cmd/sfu-server/main.go


clean:
	rm -f ./bin/*
