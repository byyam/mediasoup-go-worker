package mediasoup_go_worker

import (
	"github.com/byyam/mediasoup-go-worker/global"
	"github.com/byyam/mediasoup-go-worker/utils"
)

type SimpleWorker struct {
	WorkerBase
}

func NewSimpleWorker() *SimpleWorker {
	w := &SimpleWorker{
		WorkerBase: WorkerBase{
			pid:    global.Pid,
			logger: utils.NewLogger("worker"),
		},
	}
	return w
}

func (w *SimpleWorker) Start() {
	global.InitGlobal()
	w.logger.Info("worker[%d] start", w.pid)
}
