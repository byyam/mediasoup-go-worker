package main

import (
	"fmt"
	"github.com/byyam/mediasoup-go-worker/internal/constant"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/byyam/mediasoup-go-worker/pkg/netparser"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/google/gops/agent"
	"github.com/hashicorp/go-version"
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
	// init configurations
	mediasoupVersion := os.Getenv("MEDIASOUP_VERSION")
	currentLatest, err := version.NewVersion(mediasoupVersion)
	checkError(err)
	logger.Info("MEDIASOUP_VERSION:%s", mediasoupVersion)

	conf.InitCli()
	logger.Info("argv:%+v", conf.Settings)
	if conf.Settings.PrometheusPort > 0 {
		monitor.InitPrometheus(monitor.WithPath(conf.Settings.PrometheusPath), monitor.WithPort(conf.Settings.PrometheusPort))
	}

	// prepare write/read channel
	var netParser netparser.INetParser
	nativeJsonVersion, _ := version.NewVersion(constant.NativeJsonVersion)
	nativeVersion, _ := version.NewVersion(constant.NativeVersion)
	jsonFormat := true
	if currentLatest.GreaterThanOrEqual(nativeJsonVersion) {
		order := netparser.HostByteOrder()
		netParser, err = netparser.NewNetNativeFd(constant.ProducerChannelFd, constant.ConsumerChannelFd, order)
		logger.Info("create native codec, host order:%s", order)
		// https://github.com/versatica/mediasoup/pull/870
		if currentLatest.GreaterThanOrEqual(nativeVersion) {
			jsonFormat = false
		}
	} else {
		netParser, err = netparser.NewNetStringsFd(constant.ProducerChannelFd, constant.ConsumerChannelFd)
		logger.Info("create netstrings codec")
	}
	checkError(err)
	defer func() {
		_ = netParser.Close()
	}()

	channel := workerchannel.NewChannel(netParser, fmt.Sprintf("pid=%d,cfd=%d,pfd=%d", global.Pid, constant.ConsumerChannelFd, constant.ProducerChannelFd), jsonFormat)
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
