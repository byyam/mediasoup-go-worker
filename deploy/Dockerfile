FROM golang:1.18 as builder

ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo building for "$TARGETPLATFORM"

WORKDIR /workspace

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build -o server example/server/main.go

FROM --platform=$TARGETARCH ubuntu:latest 

WORKDIR /data/bin

COPY --from=builder /workspace/server /data/bin/server

ARG RTC_STATIC_PORT
ARG PIPE_PORT
ENV RTC_STATIC_PORT=$RTC_STATIC_PORT PIPE_PORT=$PIPE_PORT 

COPY deploy/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
