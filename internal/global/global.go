package global

import (
	"fmt"
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/utils"

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
	UdpMuxConn *ice.UDPMuxDefault
	UdpMuxPort uint16
)

func InitGlobal() {
	initUdpMuxPort()
}

func initUdpMuxPort() {
	//UdpAddr = &net.UDPAddr{Port: int(conf.Settings.RtcStaticPort)}
	if conf.Settings.RtcStaticPort != 0 { // use static port
		UdpMuxPort = conf.Settings.RtcStaticPort

		addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, UdpMuxPort)
		logger.Info("start binding static udp:%s", addr)
		if err := bindingMuxUdp(addr); err != nil {
			panic(err)
		}
	} else { // use port range
		logger.Info("start binding from port range:[%d-%d]", conf.Settings.RtcMinPort, conf.Settings.RtcMaxPort)
		for port := conf.Settings.RtcMinPort; port <= conf.Settings.RtcMaxPort; port++ {
			addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, port)
			logger.Debug("try to binding udp:%s", addr)
			if err := bindingMuxUdp(addr); err == nil {
				UdpMuxPort = port
				break
			}
			if port == conf.Settings.RtcMaxPort {
				panic("cannot binding port in range")
			}
		}
	}
	logger.Info("banding mux UDP addr:[%s:%d] success", conf.Settings.RtcListenIp, UdpMuxPort)
}

func bindingMuxUdp(addr string) (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	loggerFactory := logging.NewDefaultLoggerFactory()
	UdpMuxConn = ice.NewUDPMuxDefault(ice.UDPMuxParams{
		Logger:  loggerFactory.NewLogger("udpMux"),
		UDPConn: udpConn,
	})
	return nil
}
