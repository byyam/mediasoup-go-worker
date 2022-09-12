package main

import (
	"github.com/byyam/mediasoup-go-worker/pkg/zlog"
)

func main() {
	zlog.Init(zlog.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "example.log",
		MaxSize:               1,
		MaxBackups:            1,
		MaxAge:                1,
		LogTimeFieldFormat:    "",
		ErrorStackMarshaler:   false,
	})
	getLog := zlog.GetLogger()
	getLog.Info("this is logger")

	select {}
}
