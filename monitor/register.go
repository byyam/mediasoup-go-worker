package monitor

import "github.com/prometheus/client_golang/prometheus"

func register() {
	prometheus.MustRegister(rtpRecvCount)
	prometheus.MustRegister(rtcpRecvCount)
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

type ActionType string

const (
	ActionReceive         ActionType = "recv"
	ActionDecryptFailed              = "decrypt_failed"
	ActionUnmarshalFailed            = "unmarshal_failed"
	ActionSsrcNotFound               = "ssrc_not_found"
)
