package global

import (
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/conf"
)

var (
	Pid        = os.Getpid()
	UdpMuxConn *net.UDPConn
)

func InitGlobal() {
	var err error
	if UdpMuxConn, err = net.ListenUDP("udp", &net.UDPAddr{Port: int(conf.Settings.RtcStaticPort)}); err != nil {
		panic(err)
	}
}
