package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rtcpNamespace = "rtcp"

	rtcpRecvCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: rtcpNamespace,
		Name:      "count",
	}, []string{"action"})
)

func RtcpRecvCount(action ActionType) {
	rtcpRecvCount.WithLabelValues(string(action)).Inc()
}
