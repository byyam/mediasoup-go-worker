package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/byyam/mediasoup-go-worker/cmd/mediasoup-worker/config"
	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/google/gops/agent"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/internal/global"
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
	config.InitConfig()
	logger.Info("argv:%+v", conf.Settings)
	if conf.Settings.PrometheusPort > 0 {
		monitor.InitPrometheus(monitor.WithPath(conf.Settings.PrometheusPath), monitor.WithPort(conf.Settings.PrometheusPort))
	}

	// init worker
	mediasoupVersion := os.Getenv("MEDIASOUP_VERSION")
	channel, payloadChannel, err := mediasoup_go_worker.InitWorker(mediasoupVersion)
	checkError(err)

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
