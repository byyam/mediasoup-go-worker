package monitor

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// ice count
var (
	iceCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ice",
		Name:      "count",
	}, []string{"direction", "packet"})
)

func IceCount(direction DirectionType, packet PacketType) {
	iceCount.WithLabelValues(string(direction), string(packet)).Inc()
}

// key frame count
const (
	KeyframePkgRecv = "key_frame_recv"
	KeyframePkgSend = "key_frame_send"
	KeyframeRecvFIR = "recv_rtcp_fir"
	KeyframeRecvPLI = "recv_rtcp_pli"
	KeyframeSendFIR = "send_rtcp_fir"
	KeyframeSendPLI = "send_rtcp_pli"
)

var (
	keyframeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "keyframe",
		Name:      "count",
	}, []string{"ssrc", "event"})
)

func KeyframeCount(ssrc uint32, event string) {
	keyframeCount.WithLabelValues(fmt.Sprintf("%d", ssrc), event).Inc()
}

// mediasoup count
const (
	Router         = "router"
	Producer       = "producer"
	Consumer       = "consumer"
	SimpleConsumer = "simple_consumer"
	RtpStreamRecv  = "rtp_stream_recv"
)

var (
	mediasoupCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "mediasoup",
		Name:      "count",
	}, []string{"name", "event"})
)

func MediasoupCount(name string, event EventType) {
	mediasoupCount.WithLabelValues(name, string(event)).Inc()
}

// RTP count
var (
	rtpRecvCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rtp",
		Name:      "count",
	}, []string{"direction", "action"})
)

func RtpRecvCount(action TraceType) {
	rtpRecvCount.WithLabelValues(string(DirectionTypeRecv), string(action)).Inc()
}

func RtpSendCount(action TraceType) {
	rtpRecvCount.WithLabelValues(string(DirectionTypeSend), string(action)).Inc()
}

// RTCP count
var (
	rtcpCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rtcp",
		Name:      "count",
	}, []string{"direction", "action"})
)

func RtcpRecvCount(action TraceType) {
	rtcpCount.WithLabelValues(string(DirectionTypeRecv), string(action)).Inc()
}

func RtcpSendCount(action TraceType) {
	rtcpCount.WithLabelValues(string(DirectionTypeSend), string(action)).Inc()
}
