package workerchannel

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/byyam/mediasoup-go-worker/internal/utils"

	"github.com/ragsagar/netstringer"
)

const (
	READ_MAX_BUF = 4096
)

type Channel struct {
	consumerFd int
	producerFd int
	logger     utils.Logger

	// handle request message
	OnRequestHandler atomic.Value // func(request RequestData) ResponseData
}

func NewChannel(consumerFd, producerFd int) *Channel {
	c := &Channel{
		consumerFd: consumerFd,
		producerFd: producerFd,
		logger:     utils.NewLogger("channel"),
	}
	go c.runReadLoop()
	return c
}

func (c *Channel) writeBuf(msg []byte) (int, error) {
	return syscall.Write(c.producerFd, msg)
}

func (c *Channel) readBuf(buf []byte) (int, error) {
	return syscall.Read(c.consumerFd, buf)
}

func (c *Channel) runReadLoop() {
	buf := make([]byte, READ_MAX_BUF)
	netStringerDecoder := netstringer.NewDecoder()
	for {
		n, err := c.readBuf(buf)
		if err != nil {
			c.logger.Debug("read loop error:%v", err)
			break
		}
		netStringerDecoder.FeedData(buf[:n])
		nsPayload := <-netStringerDecoder.DataOutput
		c.processPayload(nsPayload)
	}
}

func (c *Channel) processPayload(nsPayload []byte) {
	switch nsPayload[0] {
	case '{':
		_ = c.processMessage(nsPayload)
	case 'X':
		c.logger.Debug("%s\n", nsPayload[1:])
	default:
		c.logger.Debug("unexpected data: %s", nsPayload[1:])
	}
}

func (c *Channel) OnRequest(fn func(request RequestData) ResponseData) {
	c.OnRequestHandler.Store(fn)
}

func (c *Channel) processMessage(nsPayload []byte) error {
	var reqData channelData
	if err := json.Unmarshal(nsPayload, &reqData); err != nil {
		return err
	}
	c.logger.Debug("request Id=%d, Method=%s", reqData.Id, reqData.Method)
	var ret ResponseData
	rspData := channelData{Id: reqData.Id}
	if handler, ok := c.OnRequestHandler.Load().(func(request RequestData) ResponseData); ok && handler != nil {
		var internal InternalData
		_ = internal.Unmarshal(reqData.Internal)
		ret = handler(RequestData{
			Method:   reqData.Method,
			Internal: internal,
			Data:     reqData.Data,
		})
	} else {
		rspData.Error = "OnRequestHandler not register"
	}

	if ret.Err != nil {
		c.logger.Error("response error:%v, Id=%d, Method=%s", ret.Err, reqData.Id, reqData.Method)
		rspData.Error = ret.Err.Error()
		rspData.Reason = ret.Err.Error()
	} else {
		rspData.Accepted = true
	}
	rspData.Data = ret.Data
	buf, _ := rspData.Marshal()
	n, err := c.writeBuf(buf)
	c.logger.Debug("response Id=%d, len=%d, err=%v", rspData.Id, n, err)
	return nil
}

func (c *Channel) Event(targetId int, event string) {
	msg := channelData{
		TargetId: strconv.Itoa(targetId),
		Event:    event,
	}
	jsonByte, _ := json.Marshal(&msg)
	buf := netstringer.Encode(jsonByte)
	n, err := c.writeBuf(buf)
	c.logger.Debug("send Event %+v,len=%d, err=%v", msg, n, err)
}
