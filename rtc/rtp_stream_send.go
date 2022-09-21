package rtc

import (
	"time"

	"go.uber.org/zap"

	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"

	"github.com/pion/rtcp"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

type RtpStreamSend struct {
	*RtpStream
	logger                                *zap.Logger
	lostPriorScore                        uint32 // Packets lost at last interval for score calculation.
	sentPriorScore                        uint32 // Packets sent at last interval for score calculation.
	rtxSeq                                uint16
	transmissionCounter                   *RtpDataCounter
	retransmission                        *Retransmission
	onRtpStreamRetransmitRtpPacketHandler func(packet *rtpparser.Packet)
}

type ParamRtpStreamSend struct {
	*ParamRtpStream
	bufferSize                     int
	OnRtpStreamRetransmitRtpPacket func(packet *rtpparser.Packet)
}

func newRtpStreamSend(param *ParamRtpStreamSend) *RtpStreamSend {
	r := &RtpStreamSend{
		RtpStream:                             newRtpStream(param.ParamRtpStream, 10),
		transmissionCounter:                   NewRtpDataCounter(0), // default
		retransmission:                        NewRetransmission(param.bufferSize),
		onRtpStreamRetransmitRtpPacketHandler: param.OnRtpStreamRetransmitRtpPacket,
	}
	r.logger = zaplog.NewLogger()
	return r
}

func (p *RtpStreamSend) ReceiveRtcpReceiverReport(report *rtcp.ReceptionReport) {
	nowMs := utils.GetTimeMs()
	ntp := utils.TimeMs2Ntp(nowMs)
	compactNtp := (ntp.Seconds & 0x0000FFFF) << 16
	compactNtp |= (ntp.Fractions & 0xFFFF0000) >> 16
	lastSr := report.LastSenderReport
	dlsr := report.Delay

	var rtt uint32
	if lastSr != 0 && dlsr != 0 && (compactNtp > dlsr+lastSr) {
		rtt = compactNtp - dlsr - lastSr
	}
	p.rtt = float64(rtt>>16) * 1000
	p.rtt = (float64(rtt&0x0000FFFF) / 65536) * 1000
	if p.rtt > 0 {
		p.hasRtt = true
	}
	p.packetsLost = report.TotalLost
	p.fractionLost = report.FractionLost
	p.UpdateScore(report)
}

func (p *RtpStreamSend) GetRtcpSenderReport(now time.Time) *rtcp.SenderReport {
	if p.transmissionCounter.GetPacketCount() == 0 {
		return nil
	}
	report := &rtcp.SenderReport{
		SSRC:        p.GetSsrc(),
		NTPTime:     uint64(utils.ToNtpTime(now)),
		RTPTime:     p.GetRtpTimestamp(now),
		PacketCount: uint32(p.transmissionCounter.GetPacketCount()),
		OctetCount:  uint32(p.transmissionCounter.GetBytes()),
	}
	return report
}

func (p *RtpStreamSend) GetRtcpSdesChunk() *rtcp.SourceDescription {
	return &rtcp.SourceDescription{
		Chunks: []rtcp.SourceDescriptionChunk{{
			Source: p.GetSsrc(),
			Items: []rtcp.SourceDescriptionItem{{
				Type: rtcp.SDESCNAME,
				Text: p.GetCname(),
			}},
		}},
	}
}

func (p *RtpStreamSend) ReceivePacket(packet *rtpparser.Packet) bool {
	// todo
	p.transmissionCounter.Update(packet)
	return true
}

func (p *RtpStreamSend) UpdateScore(report *rtcp.ReceptionReport) {
	// todo
	p.logger.Debug("update data by RR")
}

func (p *RtpStreamSend) ReceiveNack(nackPacket *rtcp.TransportLayerNack) {
	p.nackCount++
	for _, item := range nackPacket.Nacks {
		p.nackPacketCount += uint32(len(item.PacketList()))
		p.FillRetransmissionContainer(item)
		for _, storageItem := range p.retransmission.container {
			if storageItem == nil {
				break
			}
			packet := storageItem.packet
			p.packetsRetransmitted++
			if storageItem.sentTimes == 1 {
				p.packetsRepaired++
			}
			p.onRtpStreamRetransmitRtpPacketHandler(packet)
		}
	}
}

// This method looks for the requested RTP packets and inserts them into the
// Retransmission vector (and sets to null the next position).
//
// If RTX is used the stored packet will be RTX encoded now (if not already
// encoded in a previous resend).
func (p *RtpStreamSend) FillRetransmissionContainer(nackPair rtcp.NackPair) {
	if !p.params.UseNack {
		p.logger.Warn("NACK not supported")
		return
	}
	// Ensure the container's first element is 0.
	p.retransmission.container[0] = nil
	containerIdx := 0
	var lostSeqList []uint16
	for _, lostSeq := range nackPair.PacketList() {
		storageItem := p.retransmission.buffer[lostSeq]
		if storageItem != nil { //
			// Do nothing.
			continue
		}
		packet := storageItem.packet
		// Don't resend the packet if older than MaxRetransmissionDelay ms.
		diffTs := p.maxPacketTs - packet.Timestamp
		diffMs := diffTs * 1000 / uint32(p.params.ClockRate)
		if diffMs > MaxRetransmissionDelay {
			p.logger.Warn("ignoring retransmission for too old packet")
			continue
		}
		// Don't resent the packet if it was resent in the last RTT ms.
		nowMs := utils.GetTimeMs()
		rtt := p.rtt
		if rtt == 0 {
			rtt = DefaultRtt
		}
		if storageItem.resentAtMs != 0 && nowMs-storageItem.resentAtMs <= int64(rtt) {
			p.logger.Warn("ignoring retransmission for a packet already resent in the last RTT ms")
			continue
		}
		if p.HasRtx() {
			// todo: encode rtx
		}
		storageItem.resentAtMs = nowMs
		storageItem.sentTimes++
		p.retransmission.container[containerIdx] = storageItem
		containerIdx++
		lostSeqList = append(lostSeqList, lostSeq)
	}
	p.logger.Info("retransmission", zap.Uint16s("lost", lostSeqList), zap.Uint16s("request", nackPair.PacketList()))
}

func (p *RtpStreamSend) FillJsonStats(stat *mediasoupdata.ConsumerStat) {
	nowMs := utils.GetTimeMs()
	stat.Timestamp = nowMs
	stat.Type = "outbound-rtp"
	stat.Bitrate = p.transmissionCounter.GetBitrate(nowMs)
	stat.ByteCount = p.transmissionCounter.GetBytes()
	stat.PacketCount = p.transmissionCounter.GetPacketCount()
	stat.PacketsLost = p.packetsLost
	stat.FractionLost = p.fractionLost
	stat.NackCount = p.nackCount
	stat.NackPacketCount = p.nackPacketCount
	stat.PacketsRetransmitted = p.packetsRetransmitted
	stat.PacketsRepaired = p.packetsRepaired

	p.RtpStream.FillJsonStats(stat)
}
