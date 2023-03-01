package basehandler

import (
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/democonf"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
)

func ProducerOptions(kind mediasoupdata.MediaKind, streamId uint64, rtpParameters mediasoupdata.RtpParameters) (*mediasoupdata.ProducerOptions, error) {
	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(democonf.RouterOptions.MediaCodecs)
	if err != nil {
		logger.Error().Msgf("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, demoutils.ServerError(err)
	}
	for _, c := range routerRtpCapabilities.Codecs {
		logger.Debug().Msgf("routerRtpCapabilities:%+v", c)
	}
	rtpMapping, err := mediasoupdata.GetProducerRtpParametersMapping(
		rtpParameters, routerRtpCapabilities)
	if err != nil {
		logger.Error().Msgf("GetProducerRtpParametersMapping:%+v", err)
		return nil, demoutils.ServerError(err)
	}

	produceOptions := &mediasoupdata.ProducerOptions{
		Id:                   demoutils.GetProducerId(streamId),
		Kind:                 kind,
		RtpParameters:        rtpParameters,
		Paused:               false,
		KeyFrameRequestDelay: 0,
		AppData:              nil,
		RtpMapping:           rtpMapping,
	}
	return produceOptions, nil
}
