package rtc

import (
	"github.com/pion/rtcp"
)

type TransportCongestionControlServer struct {
	bweType                                          BweType
	maxIncomingBitrate                               uint32
	unlimitedRembCounter                             uint8
	onTransportCongestionControlServerSendRtcpPacket func(packet rtcp.Packet)
}

type TransportCongestionControlServerParam struct {
	bweType                                          BweType
	maxRtcpPacketLen                                 int64
	onTransportCongestionControlServerSendRtcpPacket func(packet rtcp.Packet)
}

func newTransportCongestionControlServer(param TransportCongestionControlServerParam) *TransportCongestionControlServer {
	return &TransportCongestionControlServer{
		bweType: param.bweType,
		onTransportCongestionControlServerSendRtcpPacket: param.onTransportCongestionControlServerSendRtcpPacket,
	}
}

func (t *TransportCongestionControlServer) GetBweType() BweType {
	return t.bweType
}

func (t *TransportCongestionControlServer) SetMaxIncomingBitrate(bitrate uint32) {

}
