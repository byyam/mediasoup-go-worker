GOARCH=
ifneq ($(arch), )
	GOARCH=$(arch)
endif

GOOS=
ifneq ($(os), )
	GOOS=$(os)
endif

MOD=mod
ifeq ($(mod),vendor)
	MOD=vendor
endif

BUILD_DATE=`date '+%Y-%m-%d_%H:%M:%S'`
GIT_HASH=`git rev-parse HEAD`
GIT_BRANCH=`git branch --show-current`
GO_VERSION=`go version|awk '{print $$3}'|sed 's/[ ][ ]*/_/g'`
LDFLAGS=-ldflags "-X=main.buildstamp=${BUILD_DATE} -X=main.githash=${GIT_HASH} -X=main.gitbranch=${GIT_BRANCH} -X=main.goversion=${GO_VERSION}"

all: mediasoup-worker sfu-server

mediasoup-worker:
	@go version
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build $(LDFLAGS) -mod=${MOD} -o bin/mediasoup-worker cmd/mediasoup-worker/main.go

sfu-server:
	@go version
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build $(LDFLAGS) -mod=${MOD} -o bin/sfu-server cmd/sfu-server/main.go


clean:
	rm -f ./bin/*
