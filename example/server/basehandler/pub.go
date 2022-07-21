package basehandler

import (
	"github.com/byyam/mediasoup-go-worker/example/internal/demoutils"
	"github.com/byyam/mediasoup-go-worker/example/server/democonf"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

func ProducerOptions(kind mediasoupdata.MediaKind, streamId uint64, rtpParameters mediasoupdata.RtpParameters) (*mediasoupdata.ProducerOptions, error) {
	routerRtpCapabilities, err := mediasoupdata.GenerateRouterRtpCapabilities(democonf.RouterOptions.MediaCodecs)
	if err != nil {
		logger.Error("GenerateRouterRtpCapabilities failed:%+v", err)
		return nil, demoutils.ServerError(err)
	}
	for _, c := range routerRtpCapabilities.Codecs {
		logger.Debug("routerRtpCapabilities:%+v", c)
	}
	rtpMapping, err := mediasoupdata.GetProducerRtpParametersMapping(
		rtpParameters, routerRtpCapabilities)
	if err != nil {
		logger.Error("GetProducerRtpParametersMapping:%+v", err)
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
