package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/byyam/mediasoup-go-worker/pkg/netparser"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/google/gops/agent"
	"github.com/hashicorp/go-version"
)

const (
	ConsumerChannelFd        = 3 // read fd
	ProducerChannelFd        = 4 // write fd
	PayloadConsumerChannelFd = 5
	PayloadProducerChannelFd = 6

	NativeVersion = "3.9.0"
)

var (
	logger = utils.NewLogger("mediasoup-worker")
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mediasoupVersion := os.Getenv("MEDIASOUP_VERSION")
	currentLatest, err := version.NewVersion(mediasoupVersion)
	checkError(err)
	logger.Info("MEDIASOUP_VERSION:%s", mediasoupVersion)

	conf.InitCli()
	logger.Info("argv:%+v", conf.Settings)
	// monitor.InitPrometheus()

	logger.Info("create producer:%d and consumer:%d socket", ProducerChannelFd, ConsumerChannelFd)

	var netParser netparser.INetParser
	nativeVersion, _ := version.NewVersion(NativeVersion)
	if currentLatest.GreaterThanOrEqual(nativeVersion) {
		order := netparser.HostByteOrder()
		netParser, err = netparser.NewNetNativeFd(ProducerChannelFd, ConsumerChannelFd, order)
		logger.Info("create native codec, host order:%s", order)
	} else {
		netParser, err = netparser.NewNetStringsFd(ProducerChannelFd, ConsumerChannelFd)
		logger.Info("create netstrings codec")
	}
	checkError(err)
	defer func() {
		_ = netParser.Close()
	}()

	channel := workerchannel.NewChannel(netParser, fmt.Sprintf("pid=%d,cfd=%d,pfd=%d", global.Pid, ConsumerChannelFd, ProducerChannelFd))
	payloadChannel := workerchannel.NewPayloadChannel()

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
