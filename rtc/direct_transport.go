package rtc

import (
	"encoding/json"

	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type DirectTransport struct {
	ITransport
	id     string
	logger zerolog.Logger
}

type directTransportParam struct {
	options mediasoupdata.DirectTransportOptions
	transportParam
}

func newDirectTransport(param directTransportParam) (ITransport, error) {
	var err error
	t := &DirectTransport{
		id:     param.Id,
		logger: zerowrapper.NewScope("direct-transport", param.Id),
	}
	param.SendRtpPacketFunc = t.SendRtpPacket
	param.SendRtcpPacketFunc = t.SendRtcpPacket
	param.SendRtcpCompoundPacketFunc = t.SendRtcpCompoundPacket
	param.NotifyCloseFunc = t.Close
	param.options.Direct = true
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	t.logger.Info().Msgf("newDirectTransport options:%# v", pretty.Formatter(param.options))
	workerchannel.RegisterHandler(param.Id, t.HandleRequest)
	return t, nil
}

func (t *DirectTransport) FillJson() json.RawMessage {
	dataDump := &FBS__Transport.DumpT{}

	t.ITransport.GetJson(dataDump)
	data, _ := json.Marshal(dataDump)

	t.logger.Debug().Str("data", string(data)).Msg("dumpData")
	return data
}

func (t *DirectTransport) SendRtpPacket(packet *rtpparser.Packet) {
	t.logger.Info().Msg("send rtp packet")
}

func (t *DirectTransport) SendRtcpPacket(packet rtcp.Packet) {
	t.logger.Info().Msg("send rtcp packet")
}

func (t *DirectTransport) SendRtcpCompoundPacket(packets []rtcp.Packet) {
	t.logger.Info().Msg("send rtcp compound packet")
}

func (t *DirectTransport) Close() {
	t.logger.Info().Msg("direct transport closed")
}

func (t *DirectTransport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug().Str("request", request.String()).Msg("handle")

	switch request.Method {

	default:
		t.ITransport.HandleRequest(request, response)
	}
}
