package taskloop

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

var taskLoopSession *session
var logger = zerowrapper.NewScope(fmt.Sprintf("taskloop"))

type session struct {
	chanTask chan sTask

	closeOnce sync.Once
	closeChan chan struct{} // state for closing

}

type sTask struct {
	fn   func() error
	done chan error
}

func (s *session) taskLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.Info().Any("r", r).Str("stack", string(debug.Stack())).Msgf("watch panic recover")
		}
		logger.Info().Msg("task loop exited")
	}()
	for {
		select {
		case <-s.closeChan:
			return
		case t := <-s.chanTask:
			err := t.fn()
			t.done <- err
		}
	}
}

func RunTask(ctx context.Context, t func() error) error {
	done := make(chan error)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case taskLoopSession.chanTask <- sTask{t, done}:
		return <-done
	}
}

func InitTaskLoopSession() {
	logger.Info().Msg("init task loop session")
	taskLoopSession = &session{
		chanTask:  make(chan sTask),
		closeChan: make(chan struct{}),
	}
	go taskLoopSession.taskLoop()
}

func CloseTaskLoopSession() error {
	if taskLoopSession.isClosed() {
		return errors.New("already closed")
	}
	taskLoopSession.closeOnce.Do(func() {
		close(taskLoopSession.closeChan)
	})
	return nil
}

func (s *session) isClosed() bool {
	select {
	case <-taskLoopSession.closeChan:
		return true
	default:
		return false
	}
}
