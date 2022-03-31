package rtc

import (
	"time"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

type RtpStreamRecv struct {
	*RtpStream
	score                            uint8
	transmissionCounter              *TransmissionCounter
	logger                           utils.Logger
	onRtpStreamSendRtcpPacketHandler func(packet rtcp.Packet)
}

type ParamRtpStreamRecv struct {
	*ParamRtpStream
	onRtpStreamSendRtcpPacket func(packet rtcp.Packet)
}

func newRtpStreamRecv(param *ParamRtpStreamRecv) *RtpStreamRecv {
	r := &RtpStreamRecv{
		RtpStream:                        newRtpStream(param.ParamRtpStream, 10),
		onRtpStreamSendRtcpPacketHandler: param.onRtpStreamSendRtcpPacket,
	}
	r.logger = utils.NewLogger("RtpStreamRecv", r.GetId())
	return r
}

func (r *RtpStreamRecv) GetScore() uint8 {
	return r.score
}

func (r *RtpStreamRecv) Pause() {}

func (r *RtpStreamRecv) ReceivePacket(packet *rtp.Packet) bool {
	if !r.RtpStream.ReceivePacket(packet) {
		r.logger.Debug("packet discarded")
		return false
	}
	// todo

	return true
}

func (r *RtpStreamRecv) ReceiveRtxPacket(packet *rtp.Packet) bool {

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

func (r *RtpStreamRecv) FillJsonStats() mediasoupdata.ProducerStat {
	return mediasoupdata.ProducerStat{
		Type:                 "inbound-rtp",
		Timestamp:            time.Now().Unix(),
		Ssrc:                 r.GetSsrc(),
		RtxSsrc:              r.GetRtxSsrc(),
		Rid:                  "",
		Kind:                 "",
		MimeType:             "",
		PacketsLost:          0,
		FractionLost:         0,
		PacketsDiscarded:     0,
		PacketsRetransmitted: 0,
		PacketsRepaired:      0,
		NackCount:            0,
		NackPacketCount:      0,
		PliCount:             0,
		FirCount:             0,
		Score:                0,
		PacketCount:          0,
		ByteCount:            0,
		Bitrate:              0,
		RoundTripTime:        0,
		RtxPacketsDiscarded:  0,
		Jitter:               0,
		BitrateByLayer:       nil,
	}
}
