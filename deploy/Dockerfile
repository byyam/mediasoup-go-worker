FROM golang:1.18 AS builder

WORKDIR /data/src

COPY . .

RUN go mod download
RUN CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -o server example/server/main.go

FROM --platform=amd64 ubuntu:latest 

WORKDIR /data/bin

COPY --from=builder /data/src/server .
COPY --from=builder /data/src/deploy/entrypoint.sh ./entrypoint.sh

RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]
