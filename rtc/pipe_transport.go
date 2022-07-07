package rtc

import (
	"encoding/json"
	"fmt"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
	"github.com/kr/pretty"
	"net"
	"strconv"
)

const (
	PipeTransportProtocol = "udp"
)

type PipeTransport struct {
	ITransport
	id     string
	logger utils.Logger

	listen mediasoupdata.TransportListenIp
	rtx    bool

	udpSocket *net.UDPConn
	udpHost   string
	udpPort   uint16
	connected *utils.AtomicBool
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
	param.NotifyCloseFunc = t.Close
	t.ITransport, err = newTransport(param.transportParam)
	if err != nil {
		return nil, err
	}
	// init pipe-transport
	t.logger.Debug("newPipeTransport options:%# v", pretty.Formatter(param.options))
	if err = t.create(&param.options); err != nil {
		return nil, err
	}
	// other options
	t.rtx = param.options.EnableRtx

	return t, nil
}

func (t *PipeTransport) create(options *mediasoupdata.PipeTransportOptions) error {
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
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		return err
	}

	t.udpPort = uint16(port)
	t.udpHost = host
	t.logger.Info("create pipe-transport addr:%s", t.udpSocket.LocalAddr())
	return nil
}

func (t *PipeTransport) Close() {
	t.logger.Warn("pipe-transport closed")
}

func (t *PipeTransport) FillJson() json.RawMessage {
	transportData := mediasoupdata.PipeTransportData{
		Tuple: mediasoupdata.TransportTuple{
			LocalIp:   t.udpHost,
			LocalPort: t.udpPort,
			Protocol:  PipeTransportProtocol,
		},
		SctpParameters: mediasoupdata.SctpParameters{},
		SctpState:      "",
		Rtx:            t.rtx,
		SrtpParameters: nil,
	}

	data, _ := json.Marshal(&transportData)
	t.logger.Debug("transportData:%+v", transportData)
	return data
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
	data := &mediasoupdata.TransportConnectData{
		Tuple: mediasoupdata.TransportTuple{
			LocalIp:    t.udpHost,
			LocalPort:  t.udpPort,
			RemoteIp:   "",
			RemotePort: 0,
			Protocol:   PipeTransportProtocol,
		},
	}

	return data, nil
}
