package rtc

import (
	"encoding/json"
	"time"

	"github.com/byyam/mediasoup-go-worker/monitor"

	"github.com/pion/rtcp"

	"github.com/pion/srtp/v2"

	"github.com/pion/rtp"

	"github.com/byyam/mediasoup-go-worker/common"

	"github.com/byyam/mediasoup-go-worker/workerchannel"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

type WebrtcTransport struct {
	ITransport
	id     string
	logger utils.Logger

	iceServer *iceServer

	dtlsTransport          *dtlsTransport
	decryptCtx, encryptCtx *srtp.Context
	connected              bool
}

type webrtcTransportParam struct {
	options mediasoupdata.WebRtcTransportOptions
	transportParam
}

func newWebrtcTransport(param webrtcTransportParam) (ITransport, error) {
	var err error
	t := &WebrtcTransport{
		id:     param.Id,
		logger: utils.NewLogger("webrtc-transport"),
	}
	param.SendRtpPacketFunc = t.SendRtpPacket
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	if t.iceServer, err = newIceServer(iceServerParam{
		iceLite:          true,
		tcp4:             true,
		OnPacketReceived: t.OnPacketReceived,
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

func (t *WebrtcTransport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug("method=%s,internal=%+v", request.Method, request.Internal)

	switch request.Method {
	case mediasoupdata.MethodTransportConnect:
		var options mediasoupdata.TransportConnectOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.Connect(options)
		response.Data, _ = json.Marshal(data)
		response.Err = err

	case mediasoupdata.MethodTransportRestartIce:

	default:
		t.ITransport.HandleRequest(request, response)
	}
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
		srtpConfig, err := t.dtlsTransport.GetSRTPConfig()
		if err != nil {
			t.logger.Error("get srtp config error:%v", err)
			return
		}
		t.decryptCtx, err = srtp.CreateContext(srtpConfig.Keys.RemoteMasterKey, srtpConfig.Keys.RemoteMasterSalt, srtpConfig.Profile)
		if err != nil {
			t.logger.Error("get srtp remote/decrypt context error:%v", err)
			return
		}
		t.encryptCtx, err = srtp.CreateContext(srtpConfig.Keys.LocalMasterKey, srtpConfig.Keys.LocalMasterSalt, srtpConfig.Profile)
		if err != nil {
			t.logger.Error("get srtp local/encrypt context error:%v", err)
			return
		}
		t.connected = true
	}()

	return t.dtlsTransport.SetRole(options.DtlsParameters)
}

func (t *WebrtcTransport) OnPacketReceived(data []byte) {
	if !t.connected {
		t.logger.Warn("webrtc not connected, ignore received packet")
		return
	}
	if utils.MatchSRTPOrSRTCP(data) {
		if !utils.IsRTCP(data) {
			monitor.RtpRecvCount(monitor.ActionReceive)
			t.OnRtpDataReceived(data) // RTP
		} else {
			monitor.RtcpRecvCount(monitor.ActionReceive)
			t.OnRtcpDataReceived(data) // RTCP
		}
	} else {
		t.logger.Warn("ignoring received packet of unknown type")
	}
}

func (t *WebrtcTransport) OnRtcpDataReceived(rawData []byte) {
	decryptHeader := &rtcp.Header{}
	decryptInput := make([]byte, len(rawData))
	actualDecrypted, err := t.decryptCtx.DecryptRTCP(decryptInput, rawData, decryptHeader)
	if err != nil {
		monitor.RtcpRecvCount(monitor.ActionDecryptFailed)
		t.logger.Error("DecryptRTCP failed:%v", err)
		return
	}

	packets, err := rtcp.Unmarshal(actualDecrypted)
	if err != nil {
		monitor.RtcpRecvCount(monitor.ActionUnmarshalFailed)
		t.logger.Error("rtcp.Unmarshal failed:%v", err)
		return
	}
	t.ITransport.ReceiveRtcpPacket(decryptHeader, packets)
}

func (t *WebrtcTransport) OnRtpDataReceived(rawData []byte) {
	decryptHeader := &rtp.Header{}
	decryptInput := make([]byte, len(rawData))
	actualDecrypted, err := t.decryptCtx.DecryptRTP(decryptInput, rawData, decryptHeader)
	if err != nil {
		monitor.RtpRecvCount(monitor.ActionDecryptFailed)
		t.logger.Error("DecryptRTP failed:%v", err)
		return
	}

	rtpPacket := &rtp.Packet{}
	if err := rtpPacket.Unmarshal(actualDecrypted); err != nil {
		monitor.RtpRecvCount(monitor.ActionUnmarshalFailed)
		t.logger.Error("rtpPacket.Unmarshal error:%v", err)
		return
	}
	t.logger.Trace("OnRtpDataReceived header%+v", rtpPacket.Header)

	t.ITransport.ReceiveRtpPacket(rtpPacket)
}

func (t *WebrtcTransport) SendRtpPacket(packet *rtp.Packet) {
	if !t.connected {
		t.logger.Warn("webrtc not connected, ignore received packet")
		return
	}
	t.logger.Trace("SendRtpPacket:%+v", packet.Header)
	decryptedRaw, err := packet.Marshal()
	if err != nil {
		t.logger.Error("rtpPacket.Marshal error:%v", err)
		return
	}
	encrypted, err := t.encryptCtx.EncryptRTP(nil, decryptedRaw, &packet.Header)
	if _, err := t.iceServer.iceConn.Write(encrypted); err != nil {
		t.logger.Error("write EncryptRTP error:%v", err)
		return
	}
}