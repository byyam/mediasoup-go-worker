package rtc

import (
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/rtpparser"
)

func ProcessRtpPacket(packet *rtpparser.Packet, mimeType mediasoupdata.RtpCodecMimeType) {
	switch mimeType.Type {
	case mediasoupdata.MimeTypeVideo:
		switch mimeType.SubType {
		case mediasoupdata.MimeSubTypeH264:
			rtpparser.ProcessRtpPacketH264(packet)
		// todo: other codecs

		default:

		}
	default:

	}
}
