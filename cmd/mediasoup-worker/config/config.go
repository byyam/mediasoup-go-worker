package config

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

func InitConfig() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "logLevel", Value: "warn", Aliases: []string{"l"}},
			&cli.StringSliceFlag{Name: "logTags", Aliases: []string{"t"}},
			&cli.StringSliceFlag{Name: "logTag"}, // mediasoup old version use logTag
			&cli.IntFlag{Name: "rtcMinPort", Value: 0, Aliases: []string{"m"}},
			&cli.IntFlag{Name: "rtcMaxPort", Value: 0, Aliases: []string{"M"}},
			&cli.StringFlag{Name: "dtlsCertificateFile", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "dtlsPrivateKeyFile", Aliases: []string{"p"}},
			&cli.IntFlag{Name: "rtcStaticPort", Value: 0, Aliases: []string{"s"}},
			&cli.IntFlag{Name: "pipePort", Value: -1, Aliases: []string{"pipeP"}},
			&cli.StringFlag{Name: "rtcListenIp", Value: "127.0.0.1", Aliases: []string{"L"}},
			&cli.StringFlag{Name: "prometheusPath", Value: "/metrics", Aliases: []string{"pm"}},
			&cli.IntFlag{Name: "prometheusPort", Value: -1, Aliases: []string{"pp"}},
		},
	}

	app.Action = func(c *cli.Context) error {
		conf.Settings.LogLevel = mediasoupdata.WorkerLogLevel(c.String("logLevel"))
		logTags := c.StringSlice("logTags")
		for _, t := range logTags {
			conf.Settings.LogTags = append(conf.Settings.LogTags, mediasoupdata.WorkerLogTag(t))
		}
		logTag := c.StringSlice("logTag")
		if len(logTag) > 0 {
			fmt.Println("use mediasoup old version: logTag option name")
			for _, t := range logTag {
				conf.Settings.LogTags = append(conf.Settings.LogTags, mediasoupdata.WorkerLogTag(t))
			}
		}
		conf.Settings.RtcMinPort = uint16(c.Int("rtcMinPort"))
		conf.Settings.RtcMaxPort = uint16(c.Int("rtcMaxPort"))
		conf.Settings.DtlsCertificateFile = c.String("dtlsCertificateFile")
		conf.Settings.DtlsPrivateKeyFile = c.String("dtlsPrivateKeyFile")
		conf.Settings.RtcStaticPort = uint16(c.Int("rtcStaticPort"))
		conf.Settings.RtcListenIp = c.String("rtcListenIp")
		conf.Settings.PrometheusPath = c.String("prometheusPath")
		conf.Settings.PrometheusPort = c.Int("prometheusPort")
		conf.Settings.PipePort = c.Int("pipePort")
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	conf.InitCli()
}
