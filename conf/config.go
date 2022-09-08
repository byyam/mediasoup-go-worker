package conf

import (
	"log"
	"sync"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

var (
	Settings mediasoupdata.WorkerSettings
	initOnce sync.Once
)

func InitCli() {
	initOnce.Do(func() {
		log.Printf("config init:%+v", Settings)
		//initLogLevel()
		//initLogScope()
		checkPort()
	})
}

//func initLogLevel() {
//	// set log level
//	switch Settings.LogLevel {
//	case mediasoupdata.WorkerLogLevel_Trace:
//		utils.ScopeLevel = utils.TraceLevel
//	case mediasoupdata.WorkerLogLevel_Debug:
//		utils.ScopeLevel = utils.DebugLevel
//	case mediasoupdata.WorkerLogLevel_Info:
//		utils.ScopeLevel = utils.InfoLevel
//	case mediasoupdata.WorkerLogLevel_Warn:
//		utils.ScopeLevel = utils.WarnLevel
//	case mediasoupdata.WorkerLogLevel_Error:
//		utils.ScopeLevel = utils.ErrorLevel
//	default:
//		panic("unknown log level")
//	}
//}
//
//func initLogScope() {
//	for _, tag := range Settings.LogTags {
//		utils.SetScopes(string(tag))
//	}
//}

func checkPort() {
	if Settings.RtcMaxPort == 0 && Settings.RtcStaticPort == 0 && Settings.RtcMinPort == 0 {
		panic("port value invalid")
	}
	if Settings.RtcMaxPort < Settings.RtcMinPort {
		panic("port range invalid")
	}
}
