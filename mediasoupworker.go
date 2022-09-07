package mediasoup_go_worker

import (
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type MediasoupWorker struct {
	workerBase
	channel        *workerchannel.Channel
	payloadChannel *workerchannel.PayloadChannel
}

func NewMediasoupWorker(channel *workerchannel.Channel, payloadChannel *workerchannel.PayloadChannel) *MediasoupWorker {
	pid := os.Getpid()
	w := &MediasoupWorker{
		workerBase: workerBase{
			pid:    pid,
			logger: utils.NewLogger("worker", pid),
		},
		channel:        channel,
		payloadChannel: payloadChannel,
	}
	w.channel.OnRequest(w.OnChannelRequest)
	return w
}

func (w *MediasoupWorker) Start() int {
	global.InitGlobal()
	w.channel.Event(w.pid, "running")
	return w.pid
}
