package rtc

import (
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

type RtpStreamRecv struct {
	*RtpStream
	score                            uint8
	logger                           utils.Logger
	onRtpStreamSendRtcpPacketHandler func(packet rtcp.Packet)
}

type ParamRtpStreamRecv struct {
	*ParamRtpStream
	onRtpStreamSendRtcpPacket func(packet rtcp.Packet)
}

func newRtpStreamRecv(param *ParamRtpStreamRecv) *RtpStreamRecv {
	r := &RtpStreamRecv{
		RtpStream:                        newRtpStream(param.ParamRtpStream),
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
		monitor.MediasoupCount(monitor.RtpStreamRecv, monitor.EventPli)
		r.onRtpStreamSendRtcpPacketHandler(packet)
	} else if r.params.UseFir {
		packet := &rtcp.FullIntraRequest{
			SenderSSRC: r.GetSsrc(),
			MediaSSRC:  r.GetSsrc(),
		}
		monitor.MediasoupCount(monitor.RtpStreamRecv, monitor.EventFir)
		r.onRtpStreamSendRtcpPacketHandler(packet)
	}
}
