package monitor

import "github.com/prometheus/client_golang/prometheus"

var (
	iceNamespace = "ice"

	iceCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: iceNamespace,
		Name:      "count",
	}, []string{"direction", "packet"})
)

func IceCount(direction DirectionType, packet PacketType) {
	iceCount.WithLabelValues(string(direction), string(packet)).Inc()
}
