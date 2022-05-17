package udpmux

import (
	"fmt"
	"net"
)

type EndPoint struct {
	remoteAddr     *net.UDPAddr
	readHandler    func(data []byte)
	onCloseHandler func()
	onWriteHandler func(data []byte, remoteAddr *net.UDPAddr) (int, error)
}

type paramEndPoint struct {
	remoteAddr string
	network    string
	onClose    func()
	onWrite    func(data []byte, remoteAddr *net.UDPAddr) (int, error)
}

func newEndPoint(param *paramEndPoint) (*EndPoint, error) {
	if param.network == "" {
		return nil, fmt.Errorf("network is empty")
	}
	remoteAddr, err := net.ResolveUDPAddr(param.network, param.remoteAddr)
	if err != nil {
		return nil, err
	}
	e := &EndPoint{
		remoteAddr:     remoteAddr,
		onCloseHandler: param.onClose,
		onWriteHandler: param.onWrite,
	}

	return e, nil
}

func (p *EndPoint) RemoteAddr() string {
	return p.remoteAddr.String()
}

func (p *EndPoint) OnRead(handler func(data []byte)) {
	p.readHandler = handler
}

func (p *EndPoint) Write(buf []byte) (int, error) {
	return p.onWriteHandler(buf, p.remoteAddr)
}

func (p *EndPoint) onRead(data []byte) {
	p.readHandler(data)
}

func (p *EndPoint) Close() {
	p.onCloseHandler()
}
