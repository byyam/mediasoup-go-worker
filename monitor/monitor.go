package monitor

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

var (
	logger = zerowrapper.NewScope("prometheus")
)

type monitorOpt struct {
	Port int
	Path string
}

func InitPrometheus(options ...func(*monitorOpt)) {
	go func() {
		register()
		// Expose the registered metrics via HTTP.
		settings := monitorOpt{
			Port: 15000,
			Path: "/metrics",
		}
		for _, option := range options {
			option(&settings)
		}
		logger.Info().Msgf("prometheus listen on http:%+v", settings)
		http.Handle(settings.Path, promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), nil))
	}()
}

type Option func(*monitorOpt)

func WithPath(path string) Option {
	return func(settings *monitorOpt) {
		if path != "" {
			settings.Path = path
		}
	}
}

func WithPort(port int) Option {
	return func(settings *monitorOpt) {
		if port != 0 {
			settings.Port = port
		}
	}
}
