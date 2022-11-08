package main

import (
	"github.com/rs/zerolog"
	"go.uber.org/zap"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/example/pkglogger/config"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger *zap.Logger
)

func main() {
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
		FileLoggingEnabled:    false,
		Directory:             "./log",
		Filename:              "signal.log",
		MaxSize:               1,
		MaxBackups:            10,
		MaxAge:                2,
		LogTimeFieldFormat:    zerolog.TimeFormatUnixMicro,
	})
	zaplog.NewLogger().Info("this is logger")
	logger = zaplog.NewLogger()
	server()

	select {}
}

func server() {
	config.InitConfig()
	worker := mediasoup_go_worker.NewSimpleWorker()
	pid := worker.Start()
	logger.Info("zap log: worker start", zap.Int("pid", pid))

	rsp := worker.OnChannelRequest(workerchannel.RequestData{
		Method:   mediasoupdata.MethodWorkerDump,
		Internal: workerchannel.InternalData{},
		Data:     nil,
	})
	logger.Error("rsp", zap.String("rsp", rsp.String()))
}
