package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rtpNamespace = "rtp"

	rtpRecvCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: rtpNamespace,
		Name:      "count",
	}, []string{"action"})
)

func RtpRecvCount(action TraceType) {
	rtpRecvCount.WithLabelValues(string(action)).Inc()
}
