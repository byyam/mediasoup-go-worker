package mediasoup_go_worker

import (
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type MediasoupWorker struct {
	workerBase
	channel        *workerchannel.Channel
	payloadChannel *workerchannel.PayloadChannel
}

func NewMediasoupWorker(channel *workerchannel.Channel, payloadChannel *workerchannel.PayloadChannel) *MediasoupWorker {
	w := &MediasoupWorker{
		workerBase: workerBase{
			pid:    global.Pid,
			logger: utils.NewLogger("worker", global.Pid),
		},
		channel:        channel,
		payloadChannel: payloadChannel,
	}
	w.channel.OnRequest(w.OnChannelRequest)
	return w
}

func (w *MediasoupWorker) Start() {
	global.InitGlobal()
	w.channel.Event(w.pid, "running")
}
