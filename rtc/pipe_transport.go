package rtc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/udpmux"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/kr/pretty"
	"github.com/pion/randutil"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/srtp/v2"
	"net"
	"strconv"
)

const (
	PipeTransportProtocol = "udp"

	srtpMasterLength      = 16
	srtpSaltLength        = 14
	srtpCryptoSuite       = srtp.ProtectionProfileAeadAes128Gcm
	srtpCryptoSuiteString = "AEAD_AES_128_GCM"
)

type PipeTransport struct {
	ITransport
	id     string
	logger utils.Logger

	listen mediasoupdata.TransportListenIp
	rtx    bool

	endpoint   *udpmux.EndPoint
	udpSocket  *net.UDPConn
	udpMuxMode bool
	connected  *utils.AtomicBool

	srtpKey                string
	srtpKeyBase64          string
	srtpSalt               string
	srtpSaltBase64         string
	decryptCtx, encryptCtx *srtp.Context

	tuple            mediasoupdata.TransportTuple
	remoteAddr       *net.UDPAddr
	remoteAddrString string
}

type pipeTransportParam struct {
	options mediasoupdata.PipeTransportOptions
	transportParam
}

func newPipeTransport(param pipeTransportParam) (ITransport, error) {
	var err error
	t := &PipeTransport{
		id:        param.Id,
		connected: &utils.AtomicBool{},
		logger:    utils.NewLogger("pipe-transport", param.Id),
	}
	param.SendRtpPacketFunc = t.SendRtpPacket
	param.SendRtcpPacketFunc = t.SendRtcpPacket
	param.SendRtcpCompoundPacketFunc = t.SendRtcpCompoundPacket
	param.NotifyCloseFunc = t.Close
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		t.logger.Error("newTransport error:%s", err.Error())
		return nil, err
	}
	// init pipe-transport
	t.logger.Info("newPipeTransport options:%# v", pretty.Formatter(param.options))
	if err = t.create(&param.options); err != nil {
		return nil, err
	}
	// other options
	t.rtx = param.options.EnableRtx
	if param.options.EnableSrtp {
		t.srtpKey, err = randutil.GenerateCryptoRandomString(srtpMasterLength, utils.RunesAlpha)
		if err != nil {
			return nil, err
		}
		t.srtpKeyBase64 = base64.StdEncoding.EncodeToString([]byte(t.srtpKey))
		t.srtpSalt, err = randutil.GenerateCryptoRandomString(srtpSaltLength, utils.RunesAlpha)
		if err != nil {
			return nil, err
		}
		t.srtpSaltBase64 = base64.StdEncoding.EncodeToString([]byte(t.srtpSalt))
	}

	return t, nil
}

func (t *PipeTransport) create(options *mediasoupdata.PipeTransportOptions) error {
	if global.UdpMuxConn == nil {
		if net.ParseIP(options.ListenIp.Ip) == nil {
			return fmt.Errorf("create pipetransport error: invalid listen ip:[%s]", options.ListenIp.Ip)
		}
		t.listen = options.ListenIp
		var addr string
		if options.Port == 0 {
			addr = fmt.Sprintf("%s:", options.ListenIp.Ip)
		} else {
			addr = fmt.Sprintf("%s:%d", options.ListenIp.Ip, options.Port)
		}
		udpAddr, err := net.ResolveUDPAddr(PipeTransportProtocol, addr)
		if err != nil {
			return err
		}
		t.udpSocket, err = net.ListenUDP(PipeTransportProtocol, udpAddr)
		host, portStr, err := net.SplitHostPort(t.udpSocket.LocalAddr().String())
		if err != nil {
			return err
		}
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return err
		}

		t.tuple.LocalPort = uint16(port)
		t.tuple.LocalIp = host
		t.logger.Info("create pipe-transport addr:[%s]", t.udpSocket.LocalAddr())
	} else {
		t.udpMuxMode = true
		t.listen = mediasoupdata.TransportListenIp{
			Ip: global.UdpMuxConn.IP(),
		}
		t.tuple.LocalIp = global.UdpMuxConn.IP()
		t.tuple.LocalPort = global.UdpMuxConn.Port()
	}
	return nil
}

func (t *PipeTransport) Close() {
	t.logger.Warn("pipe-transport closed")
}

func (t *PipeTransport) FillJson() json.RawMessage {
	transportData := mediasoupdata.PipeTransportData{
		Tuple: t.tuple,
		Rtx:   t.rtx,
	}
	// enable srtp
	if t.hasSrtp() {
		transportData.SrtpParameters = &mediasoupdata.SrtpParameters{
			CryptoSuite: srtpCryptoSuiteString,
			KeyBase64:   t.srtpKeyBase64,
			SaltBase64:  t.srtpSaltBase64,
		}
	}

	data, _ := json.Marshal(&transportData)
	t.logger.Debug("transportData:%+v", transportData)
	return data
}

func (t *PipeTransport) hasSrtp() bool {
	return t.srtpKey != ""
}

func (t *PipeTransport) HandleRequest(request workerchannel.RequestData, response *workerchannel.ResponseData) {
	t.logger.Debug("method=%s,internal=%+v", request.Method, request.Internal)

	switch request.Method {
	case mediasoupdata.MethodTransportConnect:
		var options mediasoupdata.TransportConnectOptions
		_ = json.Unmarshal(request.Data, &options)
		data, err := t.connect(options)
		response.Data, _ = json.Marshal(data)
		response.Err = err

	default:
		t.ITransport.HandleRequest(request, response)
	}
}

func (t *PipeTransport) connect(options mediasoupdata.TransportConnectOptions) (*mediasoupdata.TransportConnectData, error) {
	var err error
	if !t.hasSrtp() && options.SrtpParameters != nil {
		return nil, fmt.Errorf("connect error: invalid srtpParameters (SRTP not enabled)")
	} else if t.hasSrtp() && options.SrtpParameters == nil {
		return nil, fmt.Errorf("connect error: invalid srtpParameters (SRTP enabled)")
	} else if !t.hasSrtp() && options.SrtpParameters == nil {
		t.logger.Debug("srtp disabled")
	} else { // srtp enabled and srtp params exists
		if options.SrtpParameters.CryptoSuite != srtpCryptoSuiteString {
			return nil, fmt.Errorf("connect error: invalid/unsupported srtpParameters.cryptoSuite")
		}
		if options.SrtpParameters.KeyBase64 == "" || options.SrtpParameters.SaltBase64 == "" {
			return nil, fmt.Errorf("connect error: missing srtpParameters.keyBase64 or SaltBase64)")
		}
		srtpKey, err := base64.StdEncoding.DecodeString(options.SrtpParameters.KeyBase64)
		if err != nil {
			return nil, err
		}
		srtpSalt, err := base64.StdEncoding.DecodeString(options.SrtpParameters.SaltBase64)
		if err != nil {
			return nil, err
		}
		if len(srtpKey) != srtpMasterLength || len(srtpSalt) != srtpSaltLength {
			return nil, fmt.Errorf("connect error: invalid decoded SRTP key/salt length")
		}
		// set srtp session
		t.decryptCtx, err = srtp.CreateContext(srtpKey, srtpSalt, srtpCryptoSuite)
		if err != nil {
			t.logger.Error("get srtp remote/decrypt context error:%v", err)
			return nil, err
		}
		t.encryptCtx, err = srtp.CreateContext([]byte(t.srtpKey), []byte(t.srtpSalt), srtpCryptoSuite)
		if err != nil {
			t.logger.Error("get srtp local/encrypt context error:%v", err)
			return nil, err
		}
	}
	if !t.udpMuxMode {
		if net.ParseIP(options.Ip) == nil {
			return nil, fmt.Errorf("connect error: invalid ip:[%s]", options.Ip)
		}
		if options.Port == 0 {
			return nil, fmt.Errorf("connect error: invalid port:[%d]", options.Port)
		}
	} else {
		if t.endpoint, err = global.UdpMuxConn.AddEndPoint(options.Ip, options.Port); err != nil {
			return nil, err
		}
	}

	t.tuple.Protocol = PipeTransportProtocol
	t.tuple.RemoteIp = options.Ip
	t.tuple.RemotePort = options.Port
	t.remoteAddrString = net.JoinHostPort(options.Ip, strconv.Itoa(int(options.Port)))
	t.remoteAddr, err = net.ResolveUDPAddr(PipeTransportProtocol, t.remoteAddrString)
	t.logger.Info("pipe-transport connect addr:[%s],udpMuxMode:%v", t.remoteAddr, t.udpMuxMode)

	if !t.udpMuxMode {
		go t.udpSocketPacketReceived()
	} else {
		t.endpoint.OnRead(func(data []byte) {
			t.OnPacketReceived(data)
		})
	}

	t.ITransport.Connected()
	t.connected.Set(true)

	data := &mediasoupdata.TransportConnectData{
		Tuple: t.tuple,
	}
	return data, nil
}

func (t *PipeTransport) udpSocketPacketReceived() {
	buf := make([]byte, global.ReceiveMTU)
	for {
		n, addr, err := t.udpSocket.ReadFromUDPAddrPort(buf)
		if err != nil {
			t.logger.Warn("udpSocketPacketReceived error:%s", err.Error())
			continue
		}
		if addr.String() != t.remoteAddrString {
			t.logger.Warn("udpSocketPacketReceived error: invalid addr:[%s]", addr.String())
			continue
		}
		t.OnPacketReceived(buf[:n])
	}
}

func (t *PipeTransport) OnPacketReceived(data []byte) {
	if !t.connected.Get() {
		t.logger.Warn("pipe not connected, ignore received packet")
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

func (t *PipeTransport) SendRtpPacket(packet *rtpparser.Packet) {
	if !t.connected.Get() {
		t.logger.Warn("pipe not connected, ignore send rtp packet")
		return
	}
	t.logger.Trace("SendRtpPacket:%+v", packet.Header)
	decryptedRaw, err := packet.Marshal()
	if err != nil {
		t.logger.Error("rtpPacket.Marshal error:%v", err)
		return
	}
	if t.hasSrtp() {
		encrypted, err := t.encryptCtx.EncryptRTP(nil, decryptedRaw, &packet.Header)
		if err != nil {
			t.logger.Error("srtp encrypt error:%v", err)
			return
		}
		if _, err := t.write(encrypted); err != nil {
			t.logger.Error("write EncryptRTP error:%v", err)
			return
		}
	} else {
		if _, err := t.write(decryptedRaw); err != nil {
			t.logger.Error("write error:%v", err)
			return
		}
	}
}

func (t *PipeTransport) SendRtcpPacket(packet rtcp.Packet) {

}

func (t *PipeTransport) SendRtcpCompoundPacket(packets []rtcp.Packet) {

}

func (t *PipeTransport) write(data []byte) (int, error) {
	if t.udpMuxMode {
		return t.endpoint.Write(data)
	}
	return t.udpSocket.WriteToUDP(data, t.remoteAddr)
}

func (t *PipeTransport) OnRtpDataReceived(rawData []byte) {
	decryptHeader := &rtp.Header{}
	decryptInput := make([]byte, len(rawData))
	var rtpPacket *rtpparser.Packet
	if t.hasSrtp() {
		actualDecrypted, err := t.decryptCtx.DecryptRTP(decryptInput, rawData, decryptHeader)
		if err != nil {
			monitor.RtpRecvCount(monitor.TraceDecryptFailed)
			t.logger.Error("DecryptRTP failed:%v", err)
			return
		}
		rtpPacket, err = rtpparser.Parse(actualDecrypted)
		if err != nil {
			monitor.RtpRecvCount(monitor.TraceUnmarshalFailed)
			t.logger.Error("rtpPacket.Unmarshal error:%v", err)
			return
		}
	} else {
		var err error
		rtpPacket, err = rtpparser.Parse(rawData)
		if err != nil {
			monitor.RtpRecvCount(monitor.TraceUnmarshalFailed)
			t.logger.Error("rtpPacket.Unmarshal error:%v", err)
			return
		} // else {
		//	t.logger.Trace("rtpPacket.Unmarshal success, rtpPacket:%+v", rtpPacket.String())
		//}
	}

	t.logger.Trace("OnRtpDataReceived header%+v", rtpPacket.Header)

	t.ITransport.ReceiveRtpPacket(rtpPacket)
}

func (t *PipeTransport) OnRtcpDataReceived(rawData []byte) {
	decryptHeader := &rtcp.Header{}
	decryptInput := make([]byte, len(rawData))
	var packets []rtcp.Packet
	if t.hasSrtp() {
		actualDecrypted, err := t.decryptCtx.DecryptRTCP(decryptInput, rawData, decryptHeader)
		if err != nil {
			monitor.RtcpRecvCount(monitor.TraceDecryptFailed)
			t.logger.Error("DecryptRTCP failed:%v", err)
			return
		}
		packets, err = rtcp.Unmarshal(actualDecrypted)
		if err != nil {
			monitor.RtcpRecvCount(monitor.TraceUnmarshalFailed)
			t.logger.Error("rtcp.Unmarshal failed:%v", err)
			return
		}
	} else {
		var err error
		packets, err = rtcp.Unmarshal(rawData)
		if err != nil {
			monitor.RtcpRecvCount(monitor.TraceUnmarshalFailed)
			t.logger.Error("rtcp.Unmarshal failed:%v", err)
			return
		}
	}

	monitor.RtcpRecvCount(monitor.TraceReceive)
	t.ITransport.ReceiveRtcpPacket(decryptHeader, packets)
}
