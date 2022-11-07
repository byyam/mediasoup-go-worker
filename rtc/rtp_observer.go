package rtc

import (
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type IRtpObserver interface {
	HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData)
}

type RtpObserver struct {
	id     string
	logger zerolog.Logger
}

func (t *RtpObserver) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	defer func() {
		t.logger.Info().Str("request", request.String()).Str("response", response.String()).Msg("handle channel request done")
	}()

	switch request.Method {
	case mediasoupdata.MethodRtpObserverAddProducer:
		t.logger.Info().Msg("add producer")

	default:
		t.logger.Error().Str("method", request.Method).Msg("handle request method not found")
		return
	}
}

func newRtpObserver(id string) (IRtpObserver, error) {
	rtpObserver := &RtpObserver{
		id:     id,
		logger: zerowrapper.NewScope("rtp-observer", id),
	}
	return rtpObserver, nil
}
