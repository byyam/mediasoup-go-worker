package conf

import (
	"sync"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

var (
	Settings mediasoupdata.WorkerSettings
	initOnce sync.Once
	logger   = zerowrapper.NewScope("config")
)

func InitCli() {
	initOnce.Do(func() {
		logger.Info().Msgf("config init, settings:%+v", Settings)
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
