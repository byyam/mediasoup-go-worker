package rtc

import (
	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type AudioLevelObserver struct {
	id     string
	logger zerolog.Logger
}

type AudioLevelObserverParam struct {
	Id      string
	Options mediasoupdata.AudioLevelObserverOptions
}

func newAudioLevelObserver(param AudioLevelObserverParam) (*AudioLevelObserver, error) {
	p := &AudioLevelObserver{
		id:     param.Id,
		logger: zerowrapper.NewScope("audio-level", param.Id),
	}
	p.logger.Info().Msgf("newAudioLevelObserver options:%# v", pretty.Formatter(param.Options))
	return p, nil
}
