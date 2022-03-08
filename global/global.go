package global

import (
	"fmt"
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/conf"

	"github.com/pion/ice/v2"
	"github.com/pion/logging"
)

const (
	ReceiveMTU = 8192
)

var (
	Pid        = os.Getpid()
	UdpAddr    *net.UDPAddr
	UdpMuxConn *ice.UDPMuxDefault
	UdpConn    *net.UDPConn
)

func InitGlobal() {
	var err error
	//UdpAddr = &net.UDPAddr{Port: int(conf.Settings.RtcStaticPort)}
	UdpAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, conf.Settings.RtcStaticPort))
	UdpConn, err = net.ListenUDP("udp", UdpAddr)
	if err != nil {
		panic(err)
	}
	loggerFactory := logging.NewDefaultLoggerFactory()
	UdpMuxConn = ice.NewUDPMuxDefault(ice.UDPMuxParams{
		Logger:  loggerFactory.NewLogger("udpMux"),
		UDPConn: UdpConn,
	})
}
