package rtc

import (
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
	"github.com/pion/rtcp"
)

type RtpStreamSend struct {
	*RtpStream
	logger              utils.Logger
	bufferSize          int
	lostPriorScore      uint32 // Packets lost at last interval for score calculation.
	sentPriorScore      uint32 // Packets sent at last interval for score calculation.
	buffer              []*StorageItem
	bufferStartIdx      uint16
	storage             []StorageItem
	rtxSeq              uint16
	transmissionCounter *RtpDataCounter
}

type ParamRtpStreamSend struct {
	*ParamRtpStream
	bufferSize int
}

func newRtpStreamSend(param *ParamRtpStreamSend) *RtpStreamSend {
	r := &RtpStreamSend{
		RtpStream:           newRtpStream(param.ParamRtpStream, 10),
		transmissionCounter: NewRtpDataCounter(0), // default
	}
	r.logger = utils.NewLogger("RtpStreamSend", r.GetId())
	if param.bufferSize > 0 {
		r.bufferSize = 65536
	} else {
		r.bufferSize = 0
	}
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

func (p *RtpStreamSend) ReceivePacket(packet *rtpparser.Packet) bool {
	// todo
	p.transmissionCounter.Update(packet)
	return true
}

func (p *RtpStreamSend) UpdateScore(report *rtcp.ReceptionReport) {
	// todo
	p.logger.Debug("update data by RR:%# v", *p.RtpStream)
}

func (p *RtpStreamSend) ReceiveNack(nackPacket *rtcp.TransportLayerNack) {
	p.nackCount++
	for _, item := range nackPacket.Nacks {
		p.nackPacketCount += uint32(len(item.PacketList()))
		p.FillRetransmissionContainer(item.PacketID, uint16(item.LostPackets))
		// todo
	}
}

// This method looks for the requested RTP packets and inserts them into the
// RetransmissionContainer vector (and sets to null the next position).
//
// If RTX is used the stored packet will be RTX encoded now (if not already
// encoded in a previous resend).
func (p *RtpStreamSend) FillRetransmissionContainer(seq, bitmask uint16) {

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

	p.RtpStream.FillJsonStats(stat)
}
