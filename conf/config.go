package conf

import (
	"sync"

	"go.uber.org/zap"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
)

var (
	Settings mediasoupdata.WorkerSettings
	initOnce sync.Once
	logger   = zaplog.GetLogger().With(zap.String("scope", "config"))
)

func InitCli() {
	initOnce.Do(func() {
		logger.Info("config init", zap.Any("settings", Settings))
		checkPort()
	})
}

func checkPort() {
	if Settings.RtcMaxPort == 0 && Settings.RtcStaticPort == 0 && Settings.RtcMinPort == 0 {
		panic("port value invalid")
	}
	if Settings.RtcMaxPort < Settings.RtcMinPort {
		panic("port range invalid")
	}
}
