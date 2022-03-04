package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/byyam/mediasoup-go-worker/global"

	"github.com/byyam/mediasoup-go-worker/conf"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"

	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

const (
	ConsumerChannelFd        = 3
	ProducerChannelFd        = 4
	PayloadConsumerChannelFd = 5
	PayloadProducerChannelFd = 6
)

var (
	logger = utils.NewLogger("example")
)

func main() {
	if MediasoupVersion := os.Getenv("MEDIASOUP_VERSION"); MediasoupVersion == "" {
		panic("MEDIASOUP_VERSION incorrect")
	}

	conf.InitCli()
	logger.Debug("argv:%+v", conf.Settings)

	channel := workerchannel.NewChannel(ConsumerChannelFd, ProducerChannelFd)
	payloadChannel := workerchannel.NewPayloadChannel(PayloadConsumerChannelFd, PayloadProducerChannelFd)

	w := mediasoup_go_worker.NewWorker(channel, payloadChannel)
	w.Start()

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
