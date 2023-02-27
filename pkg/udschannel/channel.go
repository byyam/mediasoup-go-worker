package udschannel

import (
	"sync"
	"syscall"
)

type Channel struct {
	ProducerSocket int
	ConsumerSocket int

	writeLock *sync.Mutex
	readLock  *sync.Mutex
}

func NewChannel(producerSocket, consumerSocket int) *Channel {
	return &Channel{
		ProducerSocket: producerSocket,
		ConsumerSocket: consumerSocket,
		writeLock:      new(sync.Mutex),
		readLock:       new(sync.Mutex),
	}
}

func (c *Channel) Close() error {
	err := syscall.Shutdown(c.ProducerSocket, syscall.SHUT_RDWR)
	if err != nil {
		return err
	}
	err = syscall.Shutdown(c.ConsumerSocket, syscall.SHUT_RDWR)
	return err
}

func (c *Channel) Write(bytes []byte) (int, error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return syscall.Write(c.ProducerSocket, bytes)
}

func (c *Channel) Read(buf []byte) (int, error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	return syscall.Read(c.ConsumerSocket, buf)
}
