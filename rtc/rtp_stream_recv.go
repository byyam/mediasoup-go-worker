package rtc

import (
	"github.com/byyam/mediasoup-go-worker/pkg/nack"
	"math"
	"time"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/pion/rtcp"
)

type RtpStreamRecv struct {
	*RtpStream
	score              uint8
	expectedPrior      uint32 // Packets expected at last interval.
	expectedPriorScore uint32 // Packets expected at last interval for score calculation.
	receivedPrior      uint32 // Packets received at last interval.
	receivedPriorScore uint32 // Packets received at last interval for score calculation.
	lastSrTimestamp    uint32 // The middle 32 bits out of 64 in the NTP timestamp received in the most recent sender report.
	lastSrReceived     int64  // Wallclock time representing the most recent sender report arrival.
	jitter             uint32

	nackGenerator                    *nack.NackQueue
	transmissionCounter              *TransmissionCounter // Valid media + valid RTX.
	mediaTransmissionCounter         *RtpDataCounter      // Just valid media.
	logger                           utils.Logger
	onRtpStreamSendRtcpPacketHandler func(packet rtcp.Packet)
}

type ParamRtpStreamRecv struct {
	*ParamRtpStream
	onRtpStreamSendRtcpPacket func(packet rtcp.Packet)
	sendNackDelayMs           uint32
}

func newRtpStreamRecv(param *ParamRtpStreamRecv) *RtpStreamRecv {
	r := &RtpStreamRecv{
		RtpStream:                        newRtpStream(param.ParamRtpStream, 10),
		onRtpStreamSendRtcpPacketHandler: param.onRtpStreamSendRtcpPacket,
	}
	windowSize := 2500
	if param.UseDtx {
		windowSize = 6000
	}
	r.nackGenerator = nack.NewNACKQueue(&nack.ParamNackQueue{})
	r.transmissionCounter = newTransmissionCounter(param.SpatialLayers, param.TemporalLayers, windowSize)
	r.mediaTransmissionCounter = NewRtpDataCounter(0)
	r.logger = utils.NewLogger("RtpStreamRecv", r.GetId())
	r.logger.Info("new RtpStreamRecv:%# v", *param.ParamRtpStream)
	return r
}

func (r *RtpStreamRecv) GetScore() uint8 {
	return r.score
}

func (r *RtpStreamRecv) Pause() {}

func (r *RtpStreamRecv) ReceivePacket(packet *rtpparser.Packet) bool {
	if !r.RtpStream.ReceivePacket(packet) {
		r.logger.Debug("packet discarded")
		return false
	}

	// Pass the packet to the NackGenerator.
	if r.params.UseNack {
		// foundNackPkg := r.nackGenerator.ReceivePacket(packet.Packet, packet.IsKeyFrame(), false)
		r.nackGenerator.Push(packet.SequenceNumber)
		//if !r.HasRtx() && foundNackPkg {
		//	r.packetsRetransmitted++
		//	r.packetsRepaired++
		//}
	}
	// Increase transmission counter.
	r.transmissionCounter.Update(packet)
	// todo

	return true
}

func (r *RtpStreamRecv) ReceiveRtxPacket(packet *rtpparser.Packet) bool {
	if !r.params.UseNack {
		r.logger.Warn("NACK not supported")
		return false
	}
	if packet.SSRC != r.params.RtxSsrc {
		r.logger.Warn("invalid ssrc:%d on RTX packet,expect:%d", packet.SSRC, r.params.RtxSsrc)
		return false
	}
	// Check that the payload type corresponds to the one negotiated.
	if packet.PayloadType != r.params.RtxPayloadType {
		r.logger.Warn("ignoring RTX packet with invalid payload type [ssrc:%d,seq:%d,pt:%d]", packet.SSRC, packet.SequenceNumber, packet.PayloadType)
		return false
	}
	if r.HasRtx() {

	}

	return true
}

func (r *RtpStreamRecv) RequestKeyFrame() {
	if r.params.UsePli {
		packet := &rtcp.PictureLossIndication{
			SenderSSRC: r.GetSsrc(),
			MediaSSRC:  r.GetSsrc(),
		}
		monitor.KeyframeCount(r.GetSsrc(), monitor.KeyframeSendPLI)
		r.onRtpStreamSendRtcpPacketHandler(packet)
	} else if r.params.UseFir {
		packet := &rtcp.FullIntraRequest{
			SenderSSRC: r.GetSsrc(),
			MediaSSRC:  r.GetSsrc(),
		}
		monitor.KeyframeCount(r.GetSsrc(), monitor.KeyframeSendFIR)
		r.onRtpStreamSendRtcpPacketHandler(packet)
	}
}

func (r *RtpStreamRecv) FillJsonStats(stat *mediasoupdata.ProducerStat) {
	nowMs := utils.GetTimeMs()
	stat.Type = "inbound-rtp"
	stat.Timestamp = nowMs
	stat.PacketCount = r.transmissionCounter.GetPacketCount()
	stat.ByteCount = r.transmissionCounter.GetBytes()
	stat.Bitrate = r.transmissionCounter.GetBitrate(nowMs)

	r.RtpStream.FillJsonStats(stat)
}

func (r *RtpStreamRecv) ReceiveRtcpSenderReport(report *rtcp.ReceptionReport) {
	// todo
}

func (r *RtpStreamRecv) GetRtcpReceiverReport(now time.Time, worstRemoteFractionLost uint8) *rtcp.ReceptionReport {
	report := &rtcp.ReceptionReport{
		SSRC: r.GetSsrc(),
	}
	prevPacketsLost := r.packetsLost
	// Calculate Packets Expected and Lost.
	expected := r.GetExpectedPackets()
	if int64(expected) > r.transmissionCounter.GetPacketCount() {
		r.packetsLost = expected - uint32(r.mediaTransmissionCounter.GetPacketCount())
	} else {
		r.packetsLost = 0
	}
	// Calculate Fraction Lost.
	expectedInterval := expected - r.expectedPrior
	r.expectedPrior = expected

	receivedInterval := uint32(r.mediaTransmissionCounter.GetPacketCount()) - r.receivedPrior
	r.receivedPrior = uint32(r.mediaTransmissionCounter.GetPacketCount())

	lostInterval := expectedInterval - receivedInterval
	if expectedInterval == 0 || lostInterval <= 0 {
		r.fractionLost = 0
	} else {
		r.fractionLost = uint8(math.Round(float64(lostInterval<<8) / float64(expectedInterval)))
	}
	// Worst remote fraction lost is not worse than local one.
	if worstRemoteFractionLost <= r.fractionLost {
		r.reportedPacketLost += r.packetsLost - prevPacketsLost
		report.TotalLost = r.reportedPacketLost
		report.FractionLost = r.fractionLost
	} else {
		// Recalculate packetsLost.
		newLostInterval := (uint32(worstRemoteFractionLost) * expectedInterval) >> 8
		r.reportedPacketLost += newLostInterval
		report.TotalLost = r.reportedPacketLost
		report.FractionLost = r.fractionLost
	}
	// Fill the reset of the report.
	report.LastSequenceNumber = uint32(r.maxSeq) + r.cycles
	report.Jitter = r.jitter
	if r.lastSrReceived != 0 {
		// Get delay in milliseconds.
		delayMs := utils.GetTimeMs() - r.lastSrReceived
		// Express delay in units of 1/65536 seconds.
		dlsr := (delayMs / 1000) << 16
		dlsr |= (delayMs % 1000) * 65536 / 1000
		report.Delay = uint32(dlsr)
		report.LastSenderReport = r.lastSrTimestamp
	} else {
		report.Delay = 0
		report.LastSenderReport = 0
	}
	return report
}
