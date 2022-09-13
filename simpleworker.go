package mediasoup_go_worker

import (
	"os"

	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

type SimpleWorker struct {
	workerBase
}

func NewSimpleWorker() *SimpleWorker {
	pid := os.Getpid()
	w := &SimpleWorker{
		workerBase: workerBase{
			pid:    pid,
			logger: zerowrapper.NewScope("worker", pid),
		},
	}
	return w
}

func (w *SimpleWorker) Start() int {
	global.InitGlobal()
	w.logger.Info().Int("pid", w.pid).Msg("worker start")
	return w.pid
}
