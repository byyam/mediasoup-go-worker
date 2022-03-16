package global

import (
	"fmt"
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/byyam/mediasoup-go-worker/conf"

	"github.com/pion/ice/v2"
	"github.com/pion/logging"
)

const (
	ReceiveMTU = 8192
)

var (
	logger = utils.NewLogger("mediasoup-worker")
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
	addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, conf.Settings.RtcStaticPort)
	logger.Info("binding udp:%s", addr)
	UdpAddr, err = net.ResolveUDPAddr("udp", addr)
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
