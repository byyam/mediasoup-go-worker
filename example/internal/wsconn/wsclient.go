package wsconn

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/gorilla/websocket"
	"github.com/jiyeyuran/go-protoo"
)

const (
	defaultAddr = "localhost:8080"
)

type WsClientOpt struct {
	Addr        string
	Path        string
	QueryFields map[string]string
}

type WsClient struct {
	logger utils.Logger
	WsClientOpt
}

func NewWsClient(opt WsClientOpt) *WsClient {
	client := &WsClient{
		WsClientOpt: opt,
		logger:      utils.NewLogger("websocket-client"),
	}
	if opt.Addr == "" {
		client.Addr = defaultAddr
	}
	return client
}

func (c *WsClient) Conn() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.Addr, Path: c.Path}
	q := u.Query()
	for k, v := range c.QueryFields {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	c.logger.Info("connecting to %s", u.String())

	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		c.logger.Error("dial failed:%v", err)
		return nil, err
	}
	c.logger.Info("connected to %s", u.String())
	return client, nil
}

func (c *WsClient) Request(method string, req interface{}) (protoo.Message, error) {
	msg := protoo.CreateRequest(method, req)

	client, err := c.Conn()
	rsp := protoo.Message{}
	if err != nil {
		return rsp, err
	}
	defer func() { _ = client.Close() }()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				c.logger.Error("read:%v", err)
				return
			}
			c.logger.Info("recv: %s", message)
			d := json.NewDecoder(bytes.NewReader(message))
			d.UseNumber()
			if err := d.Decode(&rsp); err != nil {
				c.logger.Error("unmarshal response error%v", err)
				continue
			}
			if rsp.Id == msg.Id {
				return
			}
		}
	}()

	if err := client.WriteMessage(websocket.TextMessage, msg.Marshal()); err != nil {
		c.logger.Error("write:%v", err)
		return rsp, err
	}
	c.logger.Info("ws request:%+v", msg.String())

	for {
		select {
		case <-done:
			c.logger.Info("ws response done")
			return rsp, nil
		case <-time.After(time.Second * 10):
			c.logger.Error("ws response timeout")
			return rsp, errors.New("rsp timeout")
		}
	}
}
