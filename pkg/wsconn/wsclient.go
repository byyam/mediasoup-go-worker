package wsconn

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jiyeyuran/go-protoo"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
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
	WsClientOpt
	logger zerolog.Logger
	conn   *websocket.Conn
	sync.Mutex
}

func NewWsClient(opt WsClientOpt) (*WsClient, error) {
	client := &WsClient{
		WsClientOpt: opt,
		logger:      zerowrapper.NewScope("websocket-client"),
	}
	if opt.Addr == "" {
		client.Addr = defaultAddr
		client.logger.Warn().Str("addr", client.Addr).Msg("WsClient address not specified, use default")
	}
	var err error
	if client.conn, err = client.connect(); err != nil {
		client.logger.Error().Err(err).Msg("WsClient connect failed")
		return nil, err
	}
	return client, nil
}

func (c *WsClient) connect() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.Addr, Path: c.Path}
	q := u.Query()
	for k, v := range c.QueryFields {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	c.logger.Info().Msgf("connecting to %s", u.String())

	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		c.logger.Error().Msgf("dial failed:%v", err)
		return nil, err
	}
	c.logger.Info().Msgf("connected to %s", u.String())
	return client, nil
}

func (c *WsClient) Request(method string, req interface{}) (protoo.Message, error) {
	msg := protoo.CreateRequest(method, req)

	var rsp protoo.Message

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			c.Lock()
			_, message, err := c.conn.ReadMessage()
			c.Unlock()
			if err != nil {
				c.logger.Error().Msgf("read:%v", err)
				return
			}
			c.logger.Info().Msgf("recv: %s", message)
			d := json.NewDecoder(bytes.NewReader(message))
			d.UseNumber()
			rsp = protoo.Message{} // reset
			if err := d.Decode(&rsp); err != nil {
				c.logger.Error().Msgf("unmarshal response error%v", err)
				continue
			}
			if rsp.Id == msg.Id {
				return
			}
		}
	}()

	if err := c.conn.WriteMessage(websocket.TextMessage, msg.Marshal()); err != nil {
		c.logger.Error().Msgf("write:%v", err)
		return rsp, err
	}
	c.logger.Info().Msgf("ws request:%+v", msg.String())

	for {
		select {
		case <-done:
			c.logger.Info().Msg("ws response done")
			return rsp, nil
		case <-time.After(time.Second * 10):
			c.logger.Error().Msg("ws response timeout")
			return rsp, errors.New("rsp timeout")
		}
	}
}

func (c *WsClient) Close() error {
	return c.conn.Close()
}
