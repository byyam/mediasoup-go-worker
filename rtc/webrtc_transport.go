package rtc

import (
	"encoding/json"
	"time"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"

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

type webrtcTransportParam struct {
	options mediasoupdata.WebRtcTransportOptions
	transportParam
}

func newWebrtcTransport(id string, param webrtcTransportParam) (ITransport, error) {
	var err error
	t := &WebrtcTransport{
		id:     id,
		logger: utils.NewLogger("webrtc-transport"),
	}
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	if t.iceServer, err = newIceServer(iceServerParam{
		iceLite:          true,
		tcp4:             true,
		OnPacketReceived: t.OnRtpDataReceived,
	}); err != nil {
		return nil, err
	}
	if t.dtlsTransport, err = newDtlsTransport(dtlsTransportParam{
		role:        mediasoupdata.DtlsRole_Auto,
		connTimeout: 30 * time.Second,
	}); err != nil {
		return nil, err
	}
	t.logger.Debug("options:%+v", param.options)
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
		var options mediasoupdata.TransportConnectOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Connect(options)
		if err != nil {
			response.Data, _ = json.Marshal(data)
		}
		response.Err = err

	case mediasoupdata.MethodTransportRestartIce:

	default:
		t.ITransport.HandleRequest(request)
	}

	return
}

func (t *WebrtcTransport) Connect(options mediasoupdata.TransportConnectOptions) (*mediasoupdata.TransportConnectData, error) {
	if options.DtlsParameters == nil {
		return nil, common.ErrInvalidParam
	}
	go func() {
		iceConn, err := t.iceServer.GetConn()
		if err != nil {
			t.logger.Error("iceConn get error:%v", err)
			return
		}
		if err = t.dtlsTransport.Connect(iceConn); err != nil {
			t.logger.Error("dtls connect error:%v", err)
			return
		}
	}()

	return t.dtlsTransport.SetRole(options.DtlsParameters)
}

func (t *WebrtcTransport) OnPacketReceived(data []byte, n int) {
	if utils.MatchSRTP(data) {
		t.OnRtpDataReceived(data, n)
	}
	// todo
}

func (t *WebrtcTransport) OnRtpDataReceived(data []byte, n int) {
	rtpPacket := &rtp.Packet{}
	if err := rtpPacket.Unmarshal(data[:n]); err != nil {
		t.logger.Error("rtpPacket.Unmarshal error:%v", err)
		return
	}
	t.ITransport.ReceiveRtpPacket(rtpPacket)
}
