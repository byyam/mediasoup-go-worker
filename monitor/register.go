package monitor

import "github.com/prometheus/client_golang/prometheus"

func register() {
	prometheus.MustRegister(rtpRecvCount)
	prometheus.MustRegister(rtcpCount)
	prometheus.MustRegister(iceCount)
	prometheus.MustRegister(mediasoupCount)
	prometheus.MustRegister(keyframeCount)
}

type DirectionType string

const (
	DirectionTypeRecv DirectionType = "recv"
	DirectionTypeSend DirectionType = "send"
)

type PacketType string

const (
	PacketAll  PacketType = "all"
	PacketStun            = "stun"
	PacketDtls            = "dtls"
)

type TraceType string

const (
	TraceReceive                TraceType = "recv"
	TraceSend                             = "send"
	TraceDecryptFailed                    = "decrypt_failed"
	TraceEncryptFailed                    = "encrypt_failed"
	TraceUnmarshalFailed                  = "unmarshal_failed"
	TraceMarshalFailed                    = "marshal_failed"
	TraceSsrcNotFound                     = "ssrc_not_found"
	TraceRtpStreamRecvFailed              = "rtp_stream_recv_failed"
	TraceRtpRtxStreamRecvFailed           = "rtp_rtx_stream_recv_failed"
	TraceRtpStreamNotFound                = "rtp_stream_not_found"
	TraceUnknownRtcpType                  = "unknown_rtcp_type"
	TraceAudio                            = "audio"
	TraceVideo                            = "video"
)

type EventType string

const (
	EventSendRtp EventType = "send_rtp"
)
