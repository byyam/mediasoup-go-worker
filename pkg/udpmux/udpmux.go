package udpmux

import (
	"net"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
)

const (
	mtu = 8192
)

type UdpMux struct {
	ip        string
	port      uint16
	localAddr *net.UDPAddr
	localSock *net.UDPConn
	endPoints sync.Map
	logger    *zap.Logger
}

func NewUdpMux(network string, ip string, port uint16) (*UdpMux, error) {
	udpAddr, err := net.ResolveUDPAddr(network, net.JoinHostPort(ip, strconv.Itoa(int(port))))
	if err != nil {
		return nil, err
	}
	udpSocket, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		return nil, err
	}

	p := &UdpMux{
		ip:        ip,
		port:      port,
		localAddr: udpAddr,
		localSock: udpSocket,
		logger:    zaplog.NewLogger(),
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
			network:    p.localAddr.Network(),
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
			p.logger.Warn("udpSocketPacketReceived error", zap.Error(err))
			continue
		}
		v, ok := p.endPoints.Load(addr.String())
		if !ok {
			p.logger.Warn("udpSocketPacketReceived error", zap.String("invalid addr", addr.String()))
			continue
		}
		endpoint := v.(*EndPoint)
		endpoint.onRead(buf[:n])
	}
}

func (p *UdpMux) write(buf []byte, remoteAddr *net.UDPAddr) (int, error) {
	return p.localSock.WriteToUDP(buf, remoteAddr)
}
