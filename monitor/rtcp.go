package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rtcpNamespace = "rtcp"

	rtcpCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: rtcpNamespace,
		Name:      "count",
	}, []string{"direction", "action"})
)

func RtcpRecvCount(action TraceType) {
	rtcpCount.WithLabelValues(string(DirectionTypeRecv), string(action)).Inc()
}

func RtcpSendCount(action TraceType) {
	rtcpCount.WithLabelValues(string(DirectionTypeSend), string(action)).Inc()
}
