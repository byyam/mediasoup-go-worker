package mediasoup_go_worker

import (
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/utils"
)

type SimpleWorker struct {
	workerBase
}

func NewSimpleWorker() *SimpleWorker {
	pid := os.Getpid()
	w := &SimpleWorker{
		workerBase: workerBase{
			pid:    pid,
			logger: utils.NewLogger("worker", pid),
		},
	}
	return w
}

func (w *SimpleWorker) Start() {
	global.InitGlobal()
	w.logger.Info("worker[%d] start", w.pid)
}
