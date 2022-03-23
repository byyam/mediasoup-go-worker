package monitor

import "github.com/prometheus/client_golang/prometheus"

func register() {
	prometheus.MustRegister(rtpRecvCount)
	prometheus.MustRegister(rtcpCount)
	prometheus.MustRegister(iceCount)
	prometheus.MustRegister(mediasoupCount)
}

type DirectionType string

const (
	DirectionTypeRecv DirectionType = "recv"
	DirectionTypeSend               = "send"
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
)

type EventType string

const (
	EventRequestKeyFrame EventType = "request_key_frame"
	EventSendRtp                   = "send_rtp"
	EventPli                       = "pli"
	EventFir                       = "fir"
)
