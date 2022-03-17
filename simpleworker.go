package mediasoup_go_worker

import (
	"github.com/byyam/mediasoup-go-worker/global"
	"github.com/byyam/mediasoup-go-worker/utils"
)

type SimpleWorker struct {
	workerBase
}

func NewSimpleWorker() *SimpleWorker {
	w := &SimpleWorker{
		workerBase: workerBase{
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
