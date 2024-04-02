package rtc

import (
	"math/rand"
	"sync"
	"time"

	"github.com/pion/interceptor/pkg/twcc"
	"github.com/pion/rtcp"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/byyam/mediasoup-go-worker/pkg/seqmgr"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

const (
	PacketArrivalTimestampWindow = 500
	LimitationRembInterval       = 1500
)

type TransportCongestionControlServer struct {
	logger                                           zerolog.Logger
	bweType                                          BweType
	maxIncomingBitrate                               uint32
	unlimitedRembCounter                             uint8
	limitationRembSentAtMs                           int64
	onTransportCongestionControlServerSendRtcpPacket func(packet rtcp.Packet)
	mapPacketArrivalTimes                            sync.Map
	transportWideSeqNumberReceived                   bool
	transportCcFeedbackWideSeqNumStart               uint16
	transportCcFeedbackSenderSsrc                    uint32
	transportCcFeedbackMediaSsrc                     uint32

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

		twccRecorder: twcc.NewRecorder(rand.Uint32()),
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

}

func (t *TransportCongestionControlServer) IncomingPacket(nowMs int64, packet *rtpparser.Packet) {
	switch t.bweType {
	case TRANSPORT_CC:
		t.logger.Debug().Uint32("ssrc", packet.SSRC).Uint16("seq", packet.SequenceNumber).Msgf("IncomingPacket")
		t.twccRecorder.Record(packet.SSRC, packet.SequenceNumber, time.Since(t.startTime).Microseconds())

		//wideSeqNumber := packet.ReadTransportWideCc01()
		//if wideSeqNumber == 0 {
		//	break
		//}
		//// Only insert the packet when receiving it for the first time.
		//if _, ok := t.mapPacketArrivalTimes.Load(wideSeqNumber); ok {
		//	break
		//}
		//t.mapPacketArrivalTimes.Store(wideSeqNumber, nowMs)
		//// We may receive packets with sequence number lower than the one in
		//// previous tcc feedback, these packets may have been reported as lost
		//// previously, therefore we need to reset the start sequence num for the
		//// next tcc feedback.
		//if !t.transportWideSeqNumberReceived || seqmgr.IsSeqLowerThanUint16(wideSeqNumber, t.transportCcFeedbackWideSeqNumStart) {
		//	t.transportCcFeedbackWideSeqNumStart = wideSeqNumber
		//}
		//t.transportWideSeqNumberReceived = true
		//
		//t.MayDropOldPacketArrivalTimes(wideSeqNumber, nowMs)
		//
		//// Update the RTCP media SSRC of the ongoing Transport-CC Feedback packet.
		//t.transportCcFeedbackSenderSsrc = 0
		//t.transportCcFeedbackMediaSsrc = packet.SSRC
		//
		//t.MaySendLimitationRembFeedback(nowMs)

	case REMB:
		t.logger.Warn().Uint32("ssrc", packet.SSRC).Msg("incoming remb not handled")

	default:

	}
}

func (t *TransportCongestionControlServer) MayDropOldPacketArrivalTimes(seqNum uint16, nowMs int64) {
	// Ignore nowMs value if it's smaller than PacketArrivalTimestampWindow in
	// order to avoid negative values (should never happen) and return early if
	// the condition is met.
	if nowMs >= PacketArrivalTimestampWindow {
		var seqList []uint16
		expiryTimestamp := nowMs - PacketArrivalTimestampWindow
		t.mapPacketArrivalTimes.Range(func(key, value any) bool {
			seq := key.(uint16)
			arrivalTime := value.(int64)
			if seq != t.transportCcFeedbackWideSeqNumStart && seqmgr.IsSeqLowerThanUint16(seq, seqNum) && arrivalTime <= expiryTimestamp {
				seqList = append(seqList, seq)
			}
			return false
		})
		// clear
		for _, seq := range seqList {
			t.mapPacketArrivalTimes.Delete(seq)
		}
	}
}

func (t *TransportCongestionControlServer) MaySendLimitationRembFeedback(nowMs int64) {
	// May fix unlimitedRembCounter.
	if t.unlimitedRembCounter > 0 && t.maxIncomingBitrate != 0 {
		t.unlimitedRembCounter = 0
	}
	// In case this is the first unlimited REMB packet, send it fast.
	if ((t.bweType != TRANSPORT_CC && t.maxIncomingBitrate != 0) || t.unlimitedRembCounter > 0) &&
		(nowMs-t.limitationRembSentAtMs > LimitationRembInterval) {
		t.logger.Debug().Msgf("sending limitation RTCP REMB packet [bitrate:%d]", t.maxIncomingBitrate)
	}
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

func (t *TransportCongestionControlServer) OnTimer() {
	for {
		time.Sleep(t.interval)
		// build and send twcc
		pkts := t.twccRecorder.BuildFeedbackPacket()
		if len(pkts) == 0 {
			continue
		}
		for _, pkt := range pkts {
			t.onTransportCongestionControlServerSendRtcpPacket(pkt)
		}
	}
}
