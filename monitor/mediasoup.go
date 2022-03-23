package monitor

import "github.com/prometheus/client_golang/prometheus"

const (
	Router         = "router"
	Producer       = "producer"
	Consumer       = "consumer"
	SimpleConsumer = "simple_consumer"
	RtpStreamRecv  = "rtp_stream_recv"
)

var (
	mediasoupNamespace = "mediasoup"

	mediasoupCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: mediasoupNamespace,
		Name:      "count",
	}, []string{"name", "event"})
)

func MediasoupCount(name string, event EventType) {
	mediasoupCount.WithLabelValues(name, string(event)).Inc()
}
