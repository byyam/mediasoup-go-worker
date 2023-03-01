package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/cmd/mediasoup-worker/config"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/google/gops/agent"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
)

var (
	logger = zerowrapper.NewScope("mediasoup-worker")
	pid    int
)

var (
	githash    string
	gitbranch  string
	buildstamp string
	goversion  string
)

func printVersion() {
	log.Printf("%11s %s", "GIT_HASH:", githash)
	log.Printf("%11s %s", "GIT_BRANCH:", gitbranch)
	log.Printf("%11s %s", "BUILD_TIME:", buildstamp)
	log.Printf("%11s %s", "GO_VERSION:", goversion)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func initLog() {
	zaplog.Init(zaplog.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "media.log",
		MaxSize:               1,
		MaxBackups:            1,
		MaxAge:                1,
		LogTimeFieldFormat:    "",
		ErrorStackMarshaler:   false,
	})
	zerowrapper.InitLog(zerowrapper.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "signal.log",
		MaxSize:               1,
		MaxBackups:            10,
		MaxAge:                2,
		LogTimeFieldFormat:    zerolog.TimeFormatUnixMicro,
	})
}

func main() {
	printVersion()
	// init configurations
	config.InitConfig()
	logger.Info().Msgf("argv:%+v", conf.Settings)
	if conf.Settings.PrometheusPort > 0 {
		monitor.InitPrometheus(monitor.WithPath(conf.Settings.PrometheusPath), monitor.WithPort(conf.Settings.PrometheusPort))
	}
	initLog()

	// init worker
	mediasoupVersion := os.Getenv("MEDIASOUP_VERSION")
	channel, payloadChannel, err := mediasoup_go_worker.InitWorker(mediasoupVersion)
	checkError(err)

	w := mediasoup_go_worker.NewMediasoupWorker(channel, payloadChannel)
	pid = w.Start()
	logger.Info().Msgf("worker[%d] start", pid)

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
	logger.Warn().Msgf("[pid=%d]stop worker", pid)
}
