package rtc

import (
	"net"
	"time"
)

type iceConn struct {
	agent         *iceServer
	bytesReceived uint64
	bytesSent     uint64
	remoteAddr    net.Addr
}

func (c *iceConn) Close() error {
	//TODO implement me
	return nil
}

func (c *iceConn) LocalAddr() net.Addr {
	return c.agent.udpConn.LocalAddr()
}

func (c *iceConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *iceConn) SetDeadline(t time.Time) error {
	//TODO implement me
	return nil
}

func (c *iceConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	return nil
}

func (c *iceConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	return nil
}

func newIceConn(remoteAddr net.Addr, agent *iceServer) *iceConn {
	return &iceConn{
		remoteAddr: remoteAddr,
		agent:      agent,
	}
}

func (c *iceConn) Write(p []byte) (int, error) {
	return c.agent.udpConn.WriteTo(p, c.remoteAddr)
}

func (c *iceConn) Read(p []byte) (int, error) {
	n, err := c.agent.buffer.Read(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}
