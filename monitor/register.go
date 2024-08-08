package monitor

import "github.com/prometheus/client_golang/prometheus"

func register() {
	prometheus.MustRegister(rtpCount)
	prometheus.MustRegister(rtpTraffic)
	prometheus.MustRegister(rtcpCount)
	prometheus.MustRegister(iceCount)
	prometheus.MustRegister(mediasoupCount)
	prometheus.MustRegister(rtcpSSRCCount)
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
	TraceRtpStream              TraceType = "rtp_stream"
	TraceRtpRtxStream           TraceType = "rtp_rtx_stream"
	// rtcp type
	TraceRtcpSourceDescription TraceType = "rtcp_source_description"
	TraceRtcpGoodbye           TraceType = "rtcp_goodbye"
)

type EventType string

const (
	EventSendRtp EventType = "send_rtp"
)
