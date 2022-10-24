package mediasoup_go_worker

import (
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/global"
)

type SimpleWorker struct {
	*workerBase
}

func NewSimpleWorker() *SimpleWorker {
	pid := os.Getpid()
	w := &SimpleWorker{
		workerBase: NewWorkerBase(pid),
	}
	return w
}

func (w *SimpleWorker) Start() int {
	global.InitGlobal()
	w.logger.Info().Int("pid", w.pid).Msg("worker start")
	return w.pid
}
