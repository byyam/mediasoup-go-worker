package democonf

import (
	"github.com/byyam/mediasoup-go-worker/h264"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

var (
	enableUdp              = true
	WebrtcTransportOptions = mediasoupdata.WebRtcTransportOptions{
		ListenIps:                       nil,
		EnableUdp:                       &enableUdp,
		EnableTcp:                       false,
		PreferUdp:                       true,
		PreferTcp:                       false,
		InitialAvailableOutgoingBitrate: 600000,
		EnableSctp:                      false,
		NumSctpStreams:                  mediasoupdata.NumSctpStreams{},
		MaxSctpMessageSize:              0,
		SctpSendBufferSize:              0,
		AppData:                         nil,
	}

	RouterOptions = mediasoupdata.RouterOptions{
		MediaCodecs: []*mediasoupdata.RtpCodecCapability{
			{
				Kind:      "audio",
				MimeType:  "audio/opus",
				ClockRate: 48000,
				Channels:  2,
			},
			{
				Kind:      "video",
				MimeType:  "video/VP8",
				ClockRate: 90000,
			},
			{
				Kind:      "video",
				MimeType:  "video/H264",
				ClockRate: 90000,
				Parameters: mediasoupdata.RtpCodecSpecificParameters{
					RtpParameter: h264.RtpParameter{
						LevelAsymmetryAllowed: 1,
						PacketizationMode:     1,
					},
				},
			},
		},
	}
)
