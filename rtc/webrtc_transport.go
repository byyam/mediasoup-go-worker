package rtc

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"go.uber.org/zap"

	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/atomicbool"
	"github.com/byyam/mediasoup-go-worker/pkg/muxpkg"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/kr/pretty"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/srtp/v2"

	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type WebrtcTransport struct {
	ITransport
	id     string
	logger zerolog.Logger

	iceServer *iceServer

	dtlsTransport          *dtlsTransport
	decryptCtx, encryptCtx *srtp.Context
	connected              *atomicbool.AtomicBool
}

type webrtcTransportParam struct {
	optionsFBS *FBS__WebRtcTransport.WebRtcTransportOptionsT
	transportParam
}

func newWebrtcTransport(param webrtcTransportParam) (ITransport, error) {
	var err error
	t := &WebrtcTransport{
		id:        param.Id,
		connected: &atomicbool.AtomicBool{},
		logger:    zerowrapper.NewScope("webrtc-transport", param.Id),
	}
	param.SendRtpPacketFunc = t.SendRtpPacket
	param.SendRtcpPacketFunc = t.SendRtcpPacket
	param.SendRtcpCompoundPacketFunc = t.SendRtcpCompoundPacket
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
		role:        FBS__WebRtcTransport.DtlsRoleAUTO,
	}); err != nil {
		return nil, err
	}
	t.logger.Debug().Msgf("newWebrtcTransport options:%# v", pretty.Formatter(param.optionsFBS))
	go func() {
		<-t.iceServer.CloseChannel()
		t.logger.Warn().Msg("ice closed")
		t.Close()
		// todo emit
	}()

	workerchannel.RegisterHandler(param.Id, t.HandleRequest)
	return t, nil
}

func (t *WebrtcTransport) FillJson() json.RawMessage {
	dataDump := &FBS__Transport.DumpT{}
	t.ITransport.GetJson(dataDump)
	webrtcTransportDump := &FBS__WebRtcTransport.DumpResponseT{
		Base:             dataDump,
		IceRole:          t.iceServer.GetRole(),
		IceParameters:    t.iceServer.GetIceParameters(),
		IceCandidates:    t.iceServer.GetLocalCandidates(),
		IceState:         t.iceServer.GetState(),
		IceSelectedTuple: t.iceServer.GetSelectedTuple(),
		DtlsParameters:   t.dtlsTransport.GetDtlsParameters(),
		DtlsState:        t.dtlsTransport.GetState(),
	}
	data, _ := json.Marshal(webrtcTransportDump)

	t.logger.Debug().Str("data", string(data)).Msg("dumpData")
	return data
}

func (t *WebrtcTransport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug().Str("request", request.String()).Msg("handle")

	switch request.MethodType {
	case FBS__Request.MethodWEBRTCTRANSPORT_CONNECT:
		requestT := request.Request.Body.Value.(*FBS__WebRtcTransport.ConnectRequestT)
		data, err := t.connect(requestT.DtlsParameters)
		response.Err = err
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyWebRtcTransport_ConnectResponse,
			Value: data,
		}
		response.RspBody = rspBody

	case FBS__Request.MethodTRANSPORT_RESTART_ICE:

	case FBS__Request.MethodTRANSPORT_GET_STATS:
		response.RspBody = &FBS__Response.BodyT{
			Type:  FBS__Response.BodyWebRtcTransport_GetStatsResponse,
			Value: t.FillJsonStats(),
		}

	default:
		t.ITransport.HandleRequest(request, response)
	}
}

func (t *WebrtcTransport) FillJsonStats() *FBS__WebRtcTransport.GetStatsResponseT {
	stats := &FBS__WebRtcTransport.GetStatsResponseT{
		Base:             t.ITransport.GetBaseStats(),
		IceRole:          t.iceServer.GetRole(),
		IceState:         t.iceServer.GetState(),
		IceSelectedTuple: t.iceServer.GetSelectedTuple(),
		DtlsState:        t.dtlsTransport.GetState(),
	}
	t.logger.Debug().Any("stats", stats).Msg("FillJsonStats")
	return stats
}

func (t *WebrtcTransport) connect(options *FBS__WebRtcTransport.DtlsParametersT) (*FBS__WebRtcTransport.ConnectResponseT, error) {
	if options == nil {
		t.logger.Error().Msgf("connect failed because DtlsParametersT nil")
		return nil, mserror.ErrInvalidParam
	}
	go func() {
		t.logger.Info().Msgf("connecting")
		iceConn, err := t.iceServer.GetConn()
		if err != nil {
			t.logger.Error().Err(err).Msg("iceConn get error")
			return
		}
		if err = t.dtlsTransport.Connect(iceConn); err != nil {
			t.logger.Error().Err(err).Msg("dtls connect error")
			return
		}
		srtpConfig, err := t.dtlsTransport.GetSRTPConfig()
		if err != nil {
			t.logger.Error().Err(err).Msg("get srtp config error")
			return
		}
		t.decryptCtx, err = srtp.CreateContext(srtpConfig.Keys.RemoteMasterKey, srtpConfig.Keys.RemoteMasterSalt, srtpConfig.Profile)
		if err != nil {
			t.logger.Error().Err(err).Msg("get srtp remote/decrypt context error")
			return
		}
		t.encryptCtx, err = srtp.CreateContext(srtpConfig.Keys.LocalMasterKey, srtpConfig.Keys.LocalMasterSalt, srtpConfig.Profile)
		if err != nil {
			t.logger.Error().Err(err).Msg("get srtp local/encrypt context error")
			return
		}
		t.connected.Set(true)
	}()

	t.ITransport.Connected()

	return t.dtlsTransport.SetRole(options)
}

func (t *WebrtcTransport) OnPacketReceived(data []byte) {
	if !t.connected.Get() {
		t.logger.Warn().Msg("webrtc not connected, ignore received packet")
		return
	}
	t.ITransport.DataReceived(len(data)) // transport stats
	if muxpkg.MatchSRTPOrSRTCP(data) {
		if !muxpkg.IsRTCP(data) {
			monitor.RtpRecvCount(0, monitor.TraceReceive, len(data))
			t.OnRtpDataReceived(data) // RTP
		} else {
			monitor.RtcpRecvCount(monitor.TraceReceive)
			t.OnRtcpDataReceived(data) // RTCP
		}
	} else {
		t.logger.Warn().Msg("ignoring received packet of unknown type")
	}
}

func (t *WebrtcTransport) OnRtcpDataReceived(rawData []byte) {
	decryptHeader := &rtcp.Header{}
	decryptInput := make([]byte, len(rawData))
	actualDecrypted, err := t.decryptCtx.DecryptRTCP(decryptInput, rawData, decryptHeader)
	if err != nil {
		monitor.RtcpRecvCount(monitor.TraceDecryptFailed)
		t.logger.Error().Err(err).Msg("DecryptRTCP failed")
		return
	}

	packets, err := rtcp.Unmarshal(actualDecrypted)
	if err != nil {
		monitor.RtcpRecvCount(monitor.TraceUnmarshalFailed)
		t.logger.Error().Err(err).Msg("rtcp.Unmarshal failed")
		return
	}
	t.ITransport.ReceiveRtcpPacket(decryptHeader, packets)
}

func (t *WebrtcTransport) OnRtpDataReceived(rawData []byte) {
	decryptHeader := &rtp.Header{}
	decryptInput := make([]byte, len(rawData))
	actualDecrypted, err := t.decryptCtx.DecryptRTP(decryptInput, rawData, decryptHeader)
	if err != nil {
		monitor.RtpRecvCount(0, monitor.TraceDecryptFailed, len(rawData))
		t.logger.Error().Err(err).Msg("DecryptRTP failed")
		return
	}
	rtpPacket, err := rtpparser.Parse(actualDecrypted)
	if err != nil {
		monitor.RtpRecvCount(0, monitor.TraceUnmarshalFailed, len(rawData))
		t.logger.Error().Err(err).Msg("rtpPacket.Unmarshal error")
		return
	}
	zaplog.NewLogger().Info("WebrtcTransport: OnRtpDataReceived", zap.String("rtpPacket", rtpPacket.String()))

	t.ITransport.ReceiveRtpPacket(rtpPacket)
}

func (t *WebrtcTransport) SendRtpPacket(packet *rtpparser.Packet) {
	if !t.connected.Get() {
		t.logger.Warn().Msg("webrtc not connected, ignore send rtp packet")
		return
	}
	zaplog.NewLogger().Info("WebrtcTransport: SendRtpPacket", zap.String("rtpPacket", packet.String()))
	decryptedRaw, err := packet.Marshal()
	if err != nil {
		t.logger.Error().Err(err).Msg("rtpPacket.Marshal error")
		return
	}
	encrypted, err := t.encryptCtx.EncryptRTP(nil, decryptedRaw, &packet.Header)
	if _, err := t.iceServer.iceConn.Write(encrypted); err != nil {
		t.logger.Error().Err(err).Msg("write EncryptRTP error")
		return
	}
	t.ITransport.DataSent(len(encrypted))
}

func (t *WebrtcTransport) SendRtcpPacket(packet rtcp.Packet) {
	if !t.connected.Get() {
		t.logger.Warn().Msg("webrtc not connected, ignore send rtcp packet")
		return
	}
	t.logger.Info().Uints32("packet", packet.DestinationSSRC()).Msg("SendRtcpPacket")
	t.logger.Debug().Msgf("SendRtcpPacket:\n%+v", packet)
	decryptedRaw, err := packet.Marshal()
	if err != nil {
		t.logger.Error().Err(err).Msg("rtcpPacket.Marshal error")
		monitor.RtcpSendCount(monitor.TraceMarshalFailed)
		return
	}
	encrypted, err := t.encryptCtx.EncryptRTCP(nil, decryptedRaw, nil)
	if _, err := t.iceServer.iceConn.Write(encrypted); err != nil {
		t.logger.Error().Err(err).Msg("write EncryptRTCP error")
		monitor.RtcpSendCount(monitor.TraceEncryptFailed)
		return
	}
	t.ITransport.DataSent(len(encrypted))
	monitor.RtcpSendCount(monitor.TraceSend)
}

func (t *WebrtcTransport) SendRtcpCompoundPacket(packets []rtcp.Packet) {
	if !t.connected.Get() {
		t.logger.Warn().Msg("webrtc not connected, ignore send rtcp packet")
		return
	}
	t.logger.Info().Msgf("SendRtcpCompoundPacket:%+v", packets)
	decryptedRaw, err := rtcp.Marshal(packets)
	if err != nil {
		t.logger.Error().Err(err).Msg("rtcpPacket.Marshal error")
		monitor.RtcpSendCount(monitor.TraceMarshalFailed)
		return
	}
	encrypted, err := t.encryptCtx.EncryptRTCP(nil, decryptedRaw, nil)
	if _, err := t.iceServer.iceConn.Write(encrypted); err != nil {
		t.logger.Error().Err(err).Msg("write EncryptRTCP error")
		monitor.RtcpSendCount(monitor.TraceEncryptFailed)
		return
	}
	t.ITransport.DataSent(len(encrypted))
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
	t.logger.Info().Msg("webrtc transport closed")
}
