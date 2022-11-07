package rtc

import (
	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type AudioLevelObserver struct {
	IRtpObserver
	id     string
	logger zerolog.Logger
}

type AudioLevelObserverParam struct {
	Id      string
	Options mediasoupdata.AudioLevelObserverOptions
}

func newAudioLevelObserver(param AudioLevelObserverParam) (IRtpObserver, error) {
	var err error
	p := &AudioLevelObserver{
		id:     param.Id,
		logger: zerowrapper.NewScope("audio-level", param.Id),
	}
	p.logger.Info().Msgf("newAudioLevelObserver options:%# v", pretty.Formatter(param.Options))
	p.IRtpObserver, err = newRtpObserver(param.Id)
	if err != nil {
		return nil, err
	}
	// register handler
	workerchannel.RegisterHandler(param.Id, p.HandleRequest)
	return p, nil
}

func (t *AudioLevelObserver) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug().Str("request", request.String()).Msg("handle")

	switch request.Method {

	default:
		t.IRtpObserver.HandleRequest(request, response)
	}
}
