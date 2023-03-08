package udschannel

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

const (
	CustomerPipeStart = 3
)

// add a file reference, in case of File object gc
var socketPairFiles []*os.File

//NOTICE :the stdioMapping param is a analogy of nodejs spawn stdio map , like stdio:["ignore", "inherit", "pipe"]
func NewController(filePath string, args []string, channelNum int) (*UdsChannel, error) {
	// create socketpair * (numOf(stdio."pipe") + num *2) descriped by stdioMapping and channelNum
	procAttr := &os.ProcAttr{
		//pass all env vars
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	// 2 means fixed read and write pipe.
	pipes := make([][2]int, channelNum*2+CustomerPipeStart)
	for i := 0; i < channelNum*2+CustomerPipeStart; i++ {
		if pair, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM|syscall.FD_CLOEXEC, 0); err == nil {
			pipes[i] = pair

			// son ony use pair[1]
			// currently we just pass pipes[3:][1]
			// the os.StartProcess will do the dup2 auto
			if i >= CustomerPipeStart {
				procAttr.Files = append(procAttr.Files, os.NewFile(uintptr(pair[1]), "pipe1"))
			}
		} else {
			panic(err)
		}
	}
	socketPairFiles = append(socketPairFiles, procAttr.Files[3:]...)

	// start worker process
	son, err := os.StartProcess(filePath, append([]string{filePath}, args...), procAttr)
	if err != nil {
		panic(err)
	}

	//  new worker object, fill Child, channels
	worker := new(UdsChannel)
	worker.Child = son
	worker.ChildExit = make(chan error)

	// newChannel
	for i := 1; i <= channelNum; i++ {
		// father always use pipes[i][0], son always use pipes[i][1]
		// for father, odd pipe as producer , even pipe ad consumer
		worker.Channels = append(worker.Channels, NewChannel(pipes[i*2+1][0], pipes[i*2+2][0]))
		if err := syscall.Close(pipes[i*2+1][1]); err != nil {
			panic(err)
		}
		if err := syscall.Close(pipes[i*2+2][1]); err != nil {
			panic(err)
		}
	}

	// wait childPid, when return , write to childExit chan
	go func() {
		procStatus, err := worker.Child.Wait()
		if err != nil {
			panic("wait son error " + err.Error())
		}
		if procStatus.ExitCode() == 0 {
			worker.ChildExit <- nil
		} else {
			worker.ChildExit <- fmt.Errorf("child exit code[%d] error", procStatus.ExitCode())
		}
	}()

	return worker, nil
}

type UdsChannel struct {
	Child       *os.Process
	ChildStdOut *os.File
	ChildStdErr *os.File

	Closed bool

	Channels []*Channel

	// watch the child process status
	ChildExit chan error
}

func (w *UdsChannel) Signal(signal os.Signal) error {
	return w.Child.Signal(signal)
}

func (w *UdsChannel) StopChild() error {
	err := w.Child.Kill()
	if err != nil {
		return err
	}

	return err
}

func (w *UdsChannel) GetChannel(index int) (*Channel, error) {
	if index > len(w.Channels)-1 {
		return nil, errors.New("bad index")
	}
	return w.Channels[index], nil
}

func (w *UdsChannel) WatchChild() chan error {
	return w.ChildExit
}
