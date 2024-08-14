package democonf

import (
	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/h264"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
)

var (
	//enableUdp              = true
	//WebrtcTransportOptions = mediasoupdata.WebRtcTransportOptions{
	//	ListenIps:                       nil,
	//	EnableUdp:                       &enableUdp,
	//	EnableTcp:                       false,
	//	PreferUdp:                       true,
	//	PreferTcp:                       false,
	//	InitialAvailableOutgoingBitrate: 600000,
	//	SctpOptions: mediasoupdata.SctpOptions{
	//		EnableSctp:         false,
	//		NumSctpStreams:     mediasoupdata.NumSctpStreams{},
	//		MaxSctpMessageSize: 0,
	//		SctpSendBufferSize: 0,
	//	},
	//	AppData: nil,
	//}
	InitialAvailableOutgoingBitrate = uint32(600000)
	WebrtcTransportOptionsFBS       = FBS__WebRtcTransport.WebRtcTransportOptionsT{
		Base: &FBS__Transport.OptionsT{
			Direct:                          false,
			MaxMessageSize:                  nil,
			InitialAvailableOutgoingBitrate: &InitialAvailableOutgoingBitrate,
			EnableSctp:                      false,
			NumSctpStreams:                  nil,
			MaxSctpMessageSize:              0,
			SctpSendBufferSize:              0,
			IsDataChannel:                   false,
		},
		Listen:    nil,
		EnableUdp: true,
		EnableTcp: false,
		PreferUdp: true,
		PreferTcp: false,
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
