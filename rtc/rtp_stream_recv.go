package rtc

import (
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

type RtpStreamRecv struct {
	*RtpStream
	score                            uint8
	transmissionCounter              *TransmissionCounter // Valid media + valid RTX.
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
	windowSize := 2500
	if param.UseDtx {
		windowSize = 6000
	}
	r.transmissionCounter = newTransmissionCounter(param.SpatialLayers, param.TemporalLayers, windowSize)
	r.logger = utils.NewLogger("RtpStreamRecv", r.GetId())
	r.logger.Info("new RtpStreamRecv:%# v", *param.ParamRtpStream)
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
	// Increase transmission counter.
	r.transmissionCounter.Update(packet)
	// todo

	return true
}

func (r *RtpStreamRecv) ReceiveRtxPacket(packet *rtp.Packet) bool {
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
