package monitor

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	keyframeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "keyframe",
		Name:      "count",
	}, []string{"ssrc", "event"})
)

const (
	KeyframePkgRecv = "key_frame_recv"
	KeyframePkgSend = "key_frame_send"
	KeyframeRecvFIR = "recv_rtcp_fir"
	KeyframeRecvPLI = "recv_rtcp_pli"
	KeyframeSendFIR = "send_rtcp_fir"
	KeyframeSendPLI = "send_rtcp_pli"
)

func KeyframeCount(ssrc uint32, event string) {
	keyframeCount.WithLabelValues(fmt.Sprintf("%d", ssrc), event).Inc()
}
