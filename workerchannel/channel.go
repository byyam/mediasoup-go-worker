package workerchannel

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync/atomic"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
)

const (
	// netstring length for a 4194304 bytes payload.
	NS_MESSAGE_MAX_LEN = 4194308
	NS_PAYLOAD_MAX_LEN = 4194304
)

type Channel struct {
	logger    utils.Logger
	netParser netparser.INetParser

	// handle request message
	OnRequestHandler atomic.Value // func(request RequestData) ResponseData
}

func NewChannel(netParser netparser.INetParser, id string) *Channel {
	c := &Channel{
		netParser: netParser,
		logger:    utils.NewLogger("channel", id),
	}
	c.logger.Info("channel start")
	go c.runReadLoop()
	return c
}

func (c *Channel) runReadLoop() {
	defer c.Close()
	payload := make([]byte, NS_PAYLOAD_MAX_LEN)
	for {
		n, err := c.netParser.ReadBuffer(payload)
		if err != nil {
			c.logger.Error("Channel error:%v", err)
			break
		}
		c.logger.Debug("payload:%s", string(payload[:n]))
		c.processPayload(payload[:n])
	}
}

func (c *Channel) processPayload(nsPayload []byte) {
	switch nsPayload[0] {
	case '{':
		if err := c.processMessage(nsPayload); err != nil {
			c.logger.Error("process message failed:%s", err)
		}
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
	c.logger.Info("request Id=%d, Method=%s", reqData.Id, reqData.Method)
	var ret ResponseData
	rspData := new(channelData)
	rspData.Id = reqData.Id
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
	jsonByte, _ := json.Marshal(&rspData)

	if len(jsonByte) > NS_MESSAGE_MAX_LEN {
		return errors.New("channel response too big")
	}
	c.logger.Trace("WriteBuffer:[%s],rspData:%+v", string(jsonByte), rspData)
	if err := c.netParser.WriteBuffer(jsonByte); err != nil {
		return err
	}
	c.logger.Info("response Id=%d,err=[%v]", rspData.Id, rspData.Error)
	return nil
}

func (c *Channel) Event(targetId int, event string) {
	msg := channelData{
		TargetId: strconv.Itoa(targetId),
		Event:    event,
	}
	jsonByte, _ := json.Marshal(&msg)
	err := c.netParser.WriteBuffer(jsonByte)
	c.logger.Info("send Event msg:%+v,err=%v", msg, err)
}

func (c *Channel) Close() {
	c.logger.Info("closed")
}
