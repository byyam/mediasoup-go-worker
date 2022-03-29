package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/google/gops/agent"
)

const (
	ConsumerChannelFd        = 3
	ProducerChannelFd        = 4
	PayloadConsumerChannelFd = 5
	PayloadProducerChannelFd = 6
)

var (
	logger = utils.NewLogger("mediasoup-worker")
)

func main() {
	if MediasoupVersion := os.Getenv("MEDIASOUP_VERSION"); MediasoupVersion == "" {
		panic("MEDIASOUP_VERSION incorrect")
	}

	conf.InitCli()
	logger.Info("argv:%+v", conf.Settings)
	monitor.InitPrometheus()

	channel := workerchannel.NewChannel(ConsumerChannelFd, ProducerChannelFd)
	payloadChannel := workerchannel.NewPayloadChannel(PayloadConsumerChannelFd, PayloadProducerChannelFd)

	w := mediasoup_go_worker.NewMediasoupWorker(channel, payloadChannel)
	w.Start()

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	// block here
	listenSignal()
	w.Stop()
}

func listenSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals
	logger.Warn("[pid=%d]stop worker", global.Pid)
}
