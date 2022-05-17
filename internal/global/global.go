package global

import (
	"fmt"
	"github.com/byyam/mediasoup-go-worker/pkg/udpmux"
	logger2 "github.com/byyam/mediasoup-go-worker/utils"
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/conf"

	"github.com/pion/ice/v2"
	"github.com/pion/logging"
)

const (
	ReceiveMTU = 8192
	Network    = "udp4"
)

var (
	logger = logger2.NewLogger("mediasoup-worker")
)

var (
	Pid        = os.Getpid()
	ICEMuxConn *ice.UDPMuxDefault
	ICEMuxPort uint16

	UdpMuxConn *udpmux.UdpMux
)

func InitGlobal() {
	initICEMuxPort()
	initUdpMuxPort()
}

func initUdpMuxPort() {
	var err error
	port := conf.Settings.PipePort
	if port < 0 {
		UdpMuxConn = nil
	}
	UdpMuxConn, err = udpmux.NewUdpMux(Network, conf.Settings.RtcListenIp, uint16(port), nil)
	if err != nil {
		panic(err)
	}
	logger.Info("banding mux UDP addr:[%s:%d] success", UdpMuxConn.IP(), UdpMuxConn.Port())
}

func initICEMuxPort() {
	//UdpAddr = &net.UDPAddr{Port: int(conf.Settings.RtcStaticPort)}
	if conf.Settings.RtcStaticPort != 0 { // use static port
		ICEMuxPort = conf.Settings.RtcStaticPort

		addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, ICEMuxPort)
		logger.Info("start binding static udp:%s", addr)
		if err := bindingICEMux(addr); err != nil {
			panic(err)
		}
	} else { // use port range
		logger.Info("start binding from port range:[%d-%d]", conf.Settings.RtcMinPort, conf.Settings.RtcMaxPort)
		for port := conf.Settings.RtcMinPort; port <= conf.Settings.RtcMaxPort; port++ {
			addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, port)
			logger.Debug("try to binding udp:%s", addr)
			if err := bindingICEMux(addr); err == nil {
				ICEMuxPort = port
				break
			}
			if port == conf.Settings.RtcMaxPort {
				panic("cannot binding port in range")
			}
		}
	}
	logger.Info("banding mux ICE UDP addr:[%s:%d] success", conf.Settings.RtcListenIp, ICEMuxPort)
}

func bindingICEMux(addr string) (err error) {
	udpAddr, err := net.ResolveUDPAddr(Network, addr)
	if err != nil {
		return
	}
	udpConn, err := net.ListenUDP(Network, udpAddr)
	if err != nil {
		return
	}
	loggerFactory := logging.NewDefaultLoggerFactory()
	ICEMuxConn = ice.NewUDPMuxDefault(ice.UDPMuxParams{
		Logger:  loggerFactory.NewLogger("udpMux"),
		UDPConn: udpConn,
	})
	return nil
}
