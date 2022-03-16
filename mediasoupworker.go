package mediasoup_go_worker

import (
	"github.com/byyam/mediasoup-go-worker/global"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type MediasoupWorker struct {
	WorkerBase
	channel        *workerchannel.Channel
	payloadChannel *workerchannel.PayloadChannel
}

func NewMediasoupWorker(channel *workerchannel.Channel, payloadChannel *workerchannel.PayloadChannel) *MediasoupWorker {
	w := &MediasoupWorker{
		WorkerBase: WorkerBase{
			pid:    global.Pid,
			logger: utils.NewLogger("worker"),
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
