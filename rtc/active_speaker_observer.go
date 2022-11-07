package rtc

import (
	"github.com/kr/pretty"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ActiveSpeakerObserver struct {
	IRtpObserver
	id     string
	logger zerolog.Logger
}

type ActiveSpeakerObserverParam struct {
	Id      string
	Options mediasoupdata.ActiveSpeakerObserverOptions
}

func newActiveSpeakerObserver(param ActiveSpeakerObserverParam) (IRtpObserver, error) {
	var err error
	p := &ActiveSpeakerObserver{
		id:     param.Id,
		logger: zerowrapper.NewScope("active-speaker", param.Id),
	}
	p.logger.Info().Msgf("newActiveSpeakerObserver options:%# v", pretty.Formatter(param.Options))
	p.IRtpObserver, err = newRtpObserver(param.Id)
	if err != nil {
		return nil, err
	}
	// register handler
	workerchannel.RegisterHandler(param.Id, p.HandleRequest)
	return p, nil
}

func (t *ActiveSpeakerObserver) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug().Str("request", request.String()).Msg("handle")

	switch request.Method {

	default:
		t.IRtpObserver.HandleRequest(request, response)
	}
}
