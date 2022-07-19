package udpmux

import (
	"github.com/byyam/mediasoup-go-worker/pkg/logwrapper"
	"net"
	"strconv"
	"sync"
)

const (
	protocol = "udp"
	mtu      = 8192
)

type UdpMux struct {
	ip        string
	port      uint16
	localAddr *net.UDPAddr
	localSock *net.UDPConn
	endPoints sync.Map
	logger    logwrapper.Logger
}

func NewUdpMux(ip string, port uint16, logger logwrapper.Logger) (*UdpMux, error) {
	if logger == nil {
		logger = logwrapper.NewLogger()
	}
	udpAddr, err := net.ResolveUDPAddr(protocol, net.JoinHostPort(ip, strconv.Itoa(int(port))))
	if err != nil {
		return nil, err
	}
	udpSocket, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		return nil, err
	}

	p := &UdpMux{
		ip:        ip,
		port:      port,
		localAddr: udpAddr,
		localSock: udpSocket,
		logger:    logger,
	}
	go p.udpSocketPacketReceived()
	return p, nil
}

func (p *UdpMux) IP() string {
	return p.ip
}

func (p *UdpMux) Port() uint16 {
	return p.port
}

func (p *UdpMux) Close() error {
	return p.localSock.Close()
}

func (p *UdpMux) AddEndPoint(ip string, port uint16) (*EndPoint, error) {
	remoteAddr := net.JoinHostPort(ip, strconv.Itoa(int(port)))
	_, ok := p.endPoints.Load(remoteAddr)
	if ok {
		return nil, ErrEndPointExists
	}
	endPoint, err := newEndPoint(
		&paramEndPoint{
			remoteAddr: remoteAddr,
			onClose: func() {
				p.endPoints.Delete(remoteAddr)
			},
			onWrite: p.write,
		},
	)
	if err != nil {
		return nil, err
	}
	p.endPoints.Store(remoteAddr, endPoint)
	return endPoint, nil
}

func (p *UdpMux) udpSocketPacketReceived() {
	buf := make([]byte, mtu)
	for {
		n, addr, err := p.localSock.ReadFromUDPAddrPort(buf)
		if err != nil {
			p.logger.Warn("udpSocketPacketReceived error:%s", err.Error())
			continue
		}
		v, ok := p.endPoints.Load(addr.String())
		if !ok {
			p.logger.Warn("udpSocketPacketReceived error: invalid addr:[%s]", addr.String())
			continue
		}
		endpoint := v.(*EndPoint)
		endpoint.onRead(buf[:n])
	}
}

func (p *UdpMux) write(buf []byte, remoteAddr *net.UDPAddr) (int, error) {
	return p.localSock.WriteToUDP(buf, remoteAddr)
}
