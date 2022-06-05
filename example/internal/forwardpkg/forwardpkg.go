package forwardpkg

import (
	"log"
	"net"

	"github.com/pion/rtp"
)

type ForwardInfo struct {
	LocalAddr  *net.UDPAddr
	RemoteAddr *net.UDPAddr
	UdpConn    *net.UDPConn
}

func (f *ForwardInfo) InitConn(local, remote string) error {
	var err error
	// Create a local addr
	if f.LocalAddr, err = net.ResolveUDPAddr("udp", local); err != nil {
		return err
	}
	if f.RemoteAddr, err = net.ResolveUDPAddr("udp", remote); err != nil {
		return err
	}
	// Dial udp
	if f.UdpConn, err = net.DialUDP("udp", f.LocalAddr, f.RemoteAddr); err != nil {
		return err
	}
	log.Println("udp connection completed")
	return nil
}

func (f *ForwardInfo) Close() error {
	if f.UdpConn != nil {
		return f.UdpConn.Close()
	}
	return nil
}

func (f *ForwardInfo) Forward(rtpPacket *rtp.Packet) {
	var n int
	var err error
	b := make([]byte, 1500)
	// Marshal into original buffer with updated PayloadType
	if n, err = rtpPacket.MarshalTo(b); err != nil {
		panic(err)
	}
	// Write
	if _, err = f.UdpConn.Write(b[:n]); err != nil {
		// For this particular example, third party applications usually timeout after a short
		// amount of time during which the user doesn't have enough time to provide the answer
		// to the browser.
		// That's why, for this particular example, the user first needs to provide the answer
		// to the browser then open the third party application. Therefore we must not kill
		// the forward on "connection refused" errors
		if opError, ok := err.(*net.OpError); ok && opError.Err.Error() == "write: connection refused" {
			log.Println(opError.Err.Error())
		}
		panic(err)
	}
	log.Println("forward rtp:", rtpPacket.Header)
}
