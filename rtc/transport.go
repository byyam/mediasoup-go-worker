package rtc

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type ITransport interface {
	Close()
	FillJson() json.RawMessage
	HandleRequest(request workerchannel.RequestData) workerchannel.ResponseData
}

type Transport struct {
	id     string
	logger utils.Logger
}

func (t *Transport) Close() {
}

func (t *Transport) FillJson() json.RawMessage {
	dumpData := mediasoupdata.TransportDump{
		Id:                      t.id,
		Direct:                  false,
		ProducerIds:             nil,
		ConsumerIds:             nil,
		MapSsrcConsumerId:       nil,
		MapRtxSsrcConsumerId:    nil,
		DataProducerIds:         nil,
		DataConsumerIds:         nil,
		RecvRtpHeaderExtensions: nil,
		RtpListener:             nil,
		SctpParameters:          mediasoupdata.SctpParameters{},
		SctpState:               "",
		SctpListener:            nil,
		TraceEventTypes:         "",
		PlainTransportDump:      nil,
		WebRtcTransportDump:     nil,
	}
	data, _ := json.Marshal(&dumpData)
	t.logger.Debug("dumpData:%+v", dumpData)
	return data
}

func newTransport() ITransport {
	transport := &Transport{
		logger: utils.NewLogger("transport"),
	}
	return transport
}
func (t *Transport) HandleRequest(request workerchannel.RequestData) (response workerchannel.ResponseData) {
	t.logger.Debug("method=%s,internal=%+v", request.Method, request.InternalData)

	switch request.Method {

	case mediasoupdata.MethodTransportDump:

	case mediasoupdata.MethodTransportClose:

	case mediasoupdata.MethodTransportProduce:

	case mediasoupdata.MethodTransportConsume:

	case mediasoupdata.MethodTransportProduceData:

	case mediasoupdata.MethodTransportConsumeData:

	case mediasoupdata.MethodTransportSetMaxIncomingBitrate:

	case mediasoupdata.MethodTransportSetMaxOutgoingBitrate:

	case mediasoupdata.MethodTransportEnableTraceEvent:

	case mediasoupdata.MethodTransportGetStats:

	default:
		t.logger.Error("unknown method:%s", request.Method)
	}
	return
}
