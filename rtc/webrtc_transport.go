package rtc

import (
	"encoding/json"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/srtp/v2"
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
		logger: utils.NewLogger("webrtc-transport", param.Id),
	}
	param.SendRtpPacketFunc = t.SendRtpPacket
	param.SendRtcpPacketFunc = t.SendRtcpPacket
	param.NotifyCloseFunc = t.Close
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	if t.iceServer, err = newIceServer(iceServerParam{
		transportId:      param.Id,
		iceLite:          true,
		tcp4:             true,
		OnPacketReceived: t.OnPacketReceived,
	}); err != nil {
		return nil, err
	}
	if t.dtlsTransport, err = newDtlsTransport(dtlsTransportParam{
		transportId: param.Id,
		role:        mediasoupdata.DtlsRole_Auto,
	}); err != nil {
		return nil, err
	}
	t.logger.Debug("newWebrtcTransport options:%# v", pretty.Formatter(param.options))
	go func() {
		<-t.iceServer.CloseChannel()
		t.logger.Warn("ice closed")
		t.Close()
		// todo emit
	}()
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
		return nil, mserror.ErrInvalidParam
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
			monitor.RtpRecvCount(monitor.TraceReceive)
			t.OnRtpDataReceived(data) // RTP
		} else {
			monitor.RtcpRecvCount(monitor.TraceReceive)
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
		monitor.RtcpRecvCount(monitor.TraceDecryptFailed)
		t.logger.Error("DecryptRTCP failed:%v", err)
		return
	}

	packets, err := rtcp.Unmarshal(actualDecrypted)
	if err != nil {
		monitor.RtcpRecvCount(monitor.TraceUnmarshalFailed)
		t.logger.Error("rtcp.Unmarshal failed:%v", err)
		return
	}
	monitor.RtcpRecvCount(monitor.TraceReceive)
	t.ITransport.ReceiveRtcpPacket(decryptHeader, packets)
}

func (t *WebrtcTransport) OnRtpDataReceived(rawData []byte) {
	decryptHeader := &rtp.Header{}
	decryptInput := make([]byte, len(rawData))
	actualDecrypted, err := t.decryptCtx.DecryptRTP(decryptInput, rawData, decryptHeader)
	if err != nil {
		monitor.RtpRecvCount(monitor.TraceDecryptFailed)
		t.logger.Error("DecryptRTP failed:%v", err)
		return
	}
	rtpPacket, err := rtpparser.Parse(actualDecrypted)
	if err != nil {
		monitor.RtpRecvCount(monitor.TraceUnmarshalFailed)
		t.logger.Error("rtpPacket.Unmarshal error:%v", err)
		return
	}
	t.logger.Trace("OnRtpDataReceived header%+v", rtpPacket.Header)

	t.ITransport.ReceiveRtpPacket(rtpPacket)
}

func (t *WebrtcTransport) SendRtpPacket(packet *rtpparser.Packet) {
	if !t.connected {
		t.logger.Warn("webrtc not connected, ignore send rtp packet")
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

func (t *WebrtcTransport) SendRtcpPacket(packet rtcp.Packet) {
	if !t.connected {
		t.logger.Warn("webrtc not connected, ignore send rtcp packet")
		return
	}
	t.logger.Debug("SendRtcpPacket")
	decryptedRaw, err := packet.Marshal()
	if err != nil {
		t.logger.Error("rtcpPacket.Marshal error:%v", err)
		monitor.RtcpSendCount(monitor.TraceMarshalFailed)
		return
	}
	encrypted, err := t.encryptCtx.EncryptRTCP(nil, decryptedRaw, nil)
	if _, err := t.iceServer.iceConn.Write(encrypted); err != nil {
		t.logger.Error("write EncryptRTCP error:%v", err)
		monitor.RtcpSendCount(monitor.TraceEncryptFailed)
		return
	}
	monitor.RtcpSendCount(monitor.TraceSend)
}

func (t *WebrtcTransport) Close() {
	if t.iceServer != nil {
		t.iceServer.Disconnect()
	}
	if t.dtlsTransport != nil {
		t.dtlsTransport.Disconnect()
	}
	t.ITransport.Close()
	t.logger.Info("webrtc transport closed")
}
