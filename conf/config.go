package conf

import (
	"io"
	"os"

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
			&cli.StringSliceFlag{Name: "logTag", Aliases: []string{"t"}},
			&cli.IntFlag{Name: "rtcMinPort", Value: 40000, Aliases: []string{"m"}},
			&cli.IntFlag{Name: "rtcMaxPort", Value: 50000, Aliases: []string{"M"}},
			&cli.StringFlag{Name: "dtlsCertificateFile", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "dtlsPrivateKeyFile", Aliases: []string{"p"}},
			&cli.IntFlag{Name: "rtcStaticPort", Value: 40000, Aliases: []string{"s"}},
			&cli.StringFlag{Name: "rtcListenIp", Value: "0.0.0.0", Aliases: []string{"L"}},
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
		return nil
	}
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		panic("usage helper")
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
