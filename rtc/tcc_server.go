package rtc

import (
	"time"

	"github.com/pion/interceptor/pkg/twcc"
	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/rtctime"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

const (
	PacketArrivalTimestampWindow = 500 // twcc
	LimitationRembInterval       = 1500
	UnlimitedRembNumPackets      = 4
)

type TransportCongestionControlServer struct {
	logger                                           zerolog.Logger
	bweType                                          BweType
	maxIncomingBitrate                               uint32
	unlimitedRembCounter                             uint8
	limitationRembSentAtMs                           uint64
	onTransportCongestionControlServerSendRtcpPacket func(packet rtcp.Packet)

	// twcc
	//mapPacketArrivalTimes                            sync.Map
	//transportWideSeqNumberReceived                   bool
	//transportCcFeedbackWideSeqNumStart               uint16
	//transportCcFeedbackSenderSsrc                    uint32
	//transportCcFeedbackMediaSsrc                     uint32

	twccRecorder *twcc.Recorder
	startTime    time.Time
	interval     time.Duration
}

type TransportCongestionControlServerParam struct {
	transportId                                      string
	bweType                                          BweType
	maxRtcpPacketLen                                 int64
	onTransportCongestionControlServerSendRtcpPacket func(packet rtcp.Packet)
}

func newTransportCongestionControlServer(param TransportCongestionControlServerParam) *TransportCongestionControlServer {
	transport := &TransportCongestionControlServer{
		logger:  zerowrapper.NewScope("tcc-server", param.transportId),
		bweType: param.bweType,
		onTransportCongestionControlServerSendRtcpPacket: param.onTransportCongestionControlServerSendRtcpPacket,

		twccRecorder: twcc.NewRecorder(0), // twcc.NewRecorder(rand.Uint32()),
		startTime:    time.Now(),
		interval:     100 * time.Millisecond,
	}

	go transport.OnTimer()

	return transport
}

func (t *TransportCongestionControlServer) GetBweType() BweType {
	return t.bweType
}

func (t *TransportCongestionControlServer) SetMaxIncomingBitrate(bitrate uint32) {
	previousMaxIncomingBitrate := t.maxIncomingBitrate
	t.maxIncomingBitrate = bitrate
	t.logger.Debug().Uint32("maxIncomingBitrate", t.maxIncomingBitrate).Uint32("previousMaxIncomingBitrate", previousMaxIncomingBitrate).Msg("SetMaxIncomingBitrate")

	if previousMaxIncomingBitrate != 0 && t.maxIncomingBitrate == 0 {
		// This is to ensure that we send N REMB packets with bitrate 0 (unlimited).
		t.unlimitedRembCounter = UnlimitedRembNumPackets
		nowMs := rtctime.GetTimeMs()
		t.MaySendLimitationRembFeedback(nowMs)
	}
}

func (t *TransportCongestionControlServer) IncomingPacket(nowMs uint64, packet *rtpparser.Packet) {
	switch t.bweType {
	case TRANSPORT_CC:
		tccExt, err := packet.ReadTransportWideCc01()
		if err != nil {
			t.logger.Warn().Uint32("ssrc", packet.SSRC).Err(err).Msgf("tcc ext error")
			break
		}

		t.logger.Debug().Uint32("ssrc", packet.SSRC).Uint16("seq", tccExt.TransportSequence).Msgf("IncomingPacket")
		t.twccRecorder.Record(packet.SSRC, tccExt.TransportSequence, time.Since(t.startTime).Microseconds())

		t.MaySendLimitationRembFeedback(nowMs)

	case REMB:
		t.logger.Warn().Uint32("ssrc", packet.SSRC).Msg("incoming remb not handled")

	default:

	}
}

func (t *TransportCongestionControlServer) MaySendLimitationRembFeedback(nowMs uint64) {
	// May fix unlimitedRembCounter.
	if t.unlimitedRembCounter > 0 && t.maxIncomingBitrate != 0 {
		t.unlimitedRembCounter = 0
	}
	// In case this is the first unlimited REMB packet, send it fast.
	if ((t.bweType != REMB && t.maxIncomingBitrate != 0) || t.unlimitedRembCounter > 0) &&
		(nowMs-t.limitationRembSentAtMs > LimitationRembInterval) {
		t.logger.Debug().Msgf("sending limitation RTCP REMB packet [bitrate:%d]", t.maxIncomingBitrate)

		// No need sender and media SSRCs.
		packet := &rtcp.ReceiverEstimatedMaximumBitrate{
			SenderSSRC: 0,
			Bitrate:    float32(t.maxIncomingBitrate),
			SSRCs:      nil,
		}
		// Notify the listener.
		t.onTransportCongestionControlServerSendRtcpPacket(packet)

		t.limitationRembSentAtMs = nowMs
		if t.unlimitedRembCounter > 0 {
			t.unlimitedRembCounter--
		}
	}
}

func (t *TransportCongestionControlServer) OnTimer() {
	for {
		time.Sleep(t.interval)
		// build and send twcc
		pkts := t.twccRecorder.BuildFeedbackPacket()
		if len(pkts) == 0 {
			continue
		}
		for _, pkt := range pkts {
			t.logger.Info().Msgf("send twcc-feedback rtcp")
			t.onTransportCongestionControlServerSendRtcpPacket(pkt)
		}
	}
}
