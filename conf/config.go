package conf

import (
	"io"
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/utils"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/urfave/cli/v2"
)

var (
	Settings mediasoupdata.WorkerSettings
)

func InitCli() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "logLevel", Value: "warn", Aliases: []string{"l"}},
			&cli.StringSliceFlag{Name: "logTags", Aliases: []string{"t"}},
			&cli.IntFlag{Name: "rtcMinPort", Value: 0, Aliases: []string{"m"}},
			&cli.IntFlag{Name: "rtcMaxPort", Value: 0, Aliases: []string{"M"}},
			&cli.StringFlag{Name: "dtlsCertificateFile", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "dtlsPrivateKeyFile", Aliases: []string{"p"}},
			&cli.IntFlag{Name: "rtcStaticPort", Value: 0, Aliases: []string{"s"}},
			&cli.StringFlag{Name: "rtcListenIp", Value: "0.0.0.0", Aliases: []string{"L"}},
			&cli.StringFlag{Name: "prometheusPath", Value: "metrics", Aliases: []string{"pm"}},
			&cli.IntFlag{Name: "prometheusPort", Value: -1, Aliases: []string{"pp"}},
		},
	}

	app.Action = func(c *cli.Context) error {
		Settings.LogLevel = mediasoupdata.WorkerLogLevel(c.String("logLevel"))
		logTags := c.StringSlice("logTag")
		for _, t := range logTags {
			Settings.LogTags = append(Settings.LogTags, mediasoupdata.WorkerLogTag(t))
		}
		Settings.RtcMinPort = uint16(c.Int("rtcMinPort"))
		Settings.RtcMaxPort = uint16(c.Int("rtcMaxPort"))
		Settings.DtlsCertificateFile = c.String("dtlsCertificateFile")
		Settings.DtlsPrivateKeyFile = c.String("dtlsPrivateKeyFile")
		Settings.RtcStaticPort = uint16(c.Int("rtcStaticPort"))
		Settings.RtcListenIp = c.String("rtcListenIp")
		Settings.PrometheusPath = c.String("prometheusPath")
		Settings.PrometheusPort = c.Int("prometheusPort")
		return nil
	}
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		panic("usage helper")
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	initLogLevel()
	initLogScope()
	checkPort()
}

func initLogLevel() {
	// set log level
	switch Settings.LogLevel {
	case mediasoupdata.WorkerLogLevel_Trace:
		utils.ScopeLevel = utils.TraceLevel
	case mediasoupdata.WorkerLogLevel_Debug:
		utils.ScopeLevel = utils.DebugLevel
	case mediasoupdata.WorkerLogLevel_Info:
		utils.ScopeLevel = utils.InfoLevel
	case mediasoupdata.WorkerLogLevel_Warn:
		utils.ScopeLevel = utils.WarnLevel
	case mediasoupdata.WorkerLogLevel_Error:
		utils.ScopeLevel = utils.ErrorLevel
	default:
		panic("unknown log level")
	}
}

func initLogScope() {
	for _, tag := range Settings.LogTags {
		utils.SetScopes(string(tag))
	}
}

func checkPort() {
	if Settings.RtcMaxPort == 0 && Settings.RtcStaticPort == 0 && Settings.RtcMinPort == 0 {
		panic("port value invalid")
	}
	if Settings.RtcMaxPort < Settings.RtcMinPort {
		panic("port range invalid")
	}
}
