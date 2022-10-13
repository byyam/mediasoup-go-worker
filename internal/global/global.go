package global

import (
	"fmt"
	"net"

	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/udpmux"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/pion/ice/v2"
	"github.com/pion/logging"
)

const (
	Network = "udp4"
)

var (
	logger = zerowrapper.NewScope("mediasoup-worker")
)

var (
	ICEMuxConn *ice.UDPMuxDefault
	ICEMuxPort uint16
	ICEMuxAddr *net.UDPAddr

	UdpMuxConn *udpmux.UdpMux
)

func InitGlobal() {
	mediasoupdata.Init()
	initICEMuxPort()
	initUdpMuxPort()
}

func initUdpMuxPort() {
	var err error
	port := conf.Settings.PipePort
	if port < 0 {
		UdpMuxConn = nil
	}
	UdpMuxConn, err = udpmux.NewUdpMux(Network, conf.Settings.RtcListenIp, uint16(port))
	if err != nil {
		panic(err)
	}
	logger.Info().Str("ip", UdpMuxConn.IP()).Uint16("port", UdpMuxConn.Port()).Msg("banding mux UDP addr success")
}

func initICEMuxPort() {
	//UdpAddr = &net.UDPAddr{Port: int(conf.Settings.RtcStaticPort)}
	if conf.Settings.RtcStaticPort != 0 { // use static port
		ICEMuxPort = conf.Settings.RtcStaticPort

		addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, ICEMuxPort)
		logger.Info().Msgf("start binding static udp,addr:%s", addr)
		if err := bindingICEMux(addr); err != nil {
			panic(err)
		}
	} else { // use port range
		logger.Info().Msgf("start binding from port range:[%d-%d]", conf.Settings.RtcMinPort, conf.Settings.RtcMaxPort)
		for port := conf.Settings.RtcMinPort; port <= conf.Settings.RtcMaxPort; port++ {
			addr := fmt.Sprintf("%s:%d", conf.Settings.RtcListenIp, port)
			logger.Info().Msgf("try to binding udp,addr:%s", addr)
			if err := bindingICEMux(addr); err == nil {
				ICEMuxPort = port
				break
			}
			if port == conf.Settings.RtcMaxPort {
				panic("cannot binding port in range")
			}
		}
	}
	logger.Info().Msgf("banding mux ICE UDP addr:[%s:%d] success", conf.Settings.RtcListenIp, ICEMuxPort)
}

func bindingICEMux(addr string) (err error) {
	ICEMuxAddr, err = net.ResolveUDPAddr(Network, addr)
	if err != nil {
		return
	}
	udpConn, err := net.ListenUDP(Network, ICEMuxAddr)
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
