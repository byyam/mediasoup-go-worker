package rtc

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"
)

type WebrtcTransport struct {
	ITransport
	id     string
	logger utils.Logger

	iceServer *iceServer

	dtlsTransport *dtlsTransport
}

func newWebrtcTransport(id string, options mediasoupdata.WebRtcTransportOptions) (ITransport, error) {
	t := &WebrtcTransport{
		ITransport: newTransport(),
		id:         id,
		logger:     utils.NewLogger("webrtc-transport"),
	}
	var err error
	if t.iceServer, err = newIceServer(iceServerParam{
		iceLite: true,
		udp4:    true,
	}); err != nil {
		return nil, err
	}
	if t.dtlsTransport, err = newDtlsTransport(dtlsTransportParam{
		isClient: false,
	}); err != nil {
		return nil, err
	}
	t.logger.Debug("options:%+v", options)
	return t, nil
}

func (t *WebrtcTransport) FillJson() json.RawMessage {
	transportData := mediasoupdata.WebrtcTransportData{
		IceRole:          t.iceServer.GetRole(),
		IceParameters:    t.iceServer.GetIceParameters(),
		IceCandidates:    t.iceServer.GetLocalCandidates(),
		IceState:         t.iceServer.GetState(),
		IceSelectedTuple: t.iceServer.GetSelectedTuple(),
		DtlsParameters:   t.dtlsTransport.GetDtlsParameters(),
		DtlsState:        t.dtlsTransport.GetState(),
		DtlsRemoteCert:   "",
		SctpParameters:   mediasoupdata.SctpParameters{},
		SctpState:        "",
	}
	data, _ := json.Marshal(&transportData)
	t.logger.Debug("transportData:%+v", transportData)
	return data
}

func (t *WebrtcTransport) HandleRequest(request workerchannel.RequestData) (response workerchannel.ResponseData) {
	t.logger.Debug("method=%s,internal=%+v", request.Method, request.InternalData)

	switch request.Method {
	case mediasoupdata.MethodTransportConnect:

	case mediasoupdata.MethodTransportRestartIce:

	default:
		t.ITransport.HandleRequest(request)
	}

	return
}
