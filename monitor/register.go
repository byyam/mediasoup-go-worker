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
	PacketStun PacketType = "stun"
	PacketDtls PacketType = "dtls"
)

type TraceType string

const (
	TraceReceive                TraceType = "recv"
	TraceSend                   TraceType = "send"
	TraceDecryptFailed          TraceType = "decrypt_failed"
	TraceEncryptFailed          TraceType = "encrypt_failed"
	TraceUnmarshalFailed        TraceType = "unmarshal_failed"
	TraceMarshalFailed          TraceType = "marshal_failed"
	TraceSsrcNotFound           TraceType = "ssrc_not_found"
	TraceRtpStreamRecvFailed    TraceType = "rtp_stream_recv_failed"
	TraceRtpRtxStreamRecvFailed TraceType = "rtp_rtx_stream_recv_failed"
	TraceRtpStreamNotFound      TraceType = "rtp_stream_not_found"
	TraceUnknownRtcpType        TraceType = "unknown_rtcp_type"
	TraceAudio                  TraceType = "audio"
	TraceVideo                  TraceType = "video"
)

type EventType string

const (
	EventSendRtp EventType = "send_rtp"
)
