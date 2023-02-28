package wsconn

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/jiyeyuran/go-protoo"
	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	logger zerolog.Logger
	WsServerOpt
}

type WsServerOpt struct {
	PingInterval   time.Duration
	PongWait       time.Duration
	Conn           *websocket.Conn
	RequestHandler func(message protoo.Message) *protoo.Message
}

func (w WsServerOpt) valid() bool {
	if w.RequestHandler == nil {
		return false
	}
	return true
}

func NewWsServer(opt WsServerOpt) (*WsServer, error) {
	if !opt.valid() {
		return nil, errors.New("input param invalid")
	}
	w := &WsServer{
		WsServerOpt: opt,
		logger:      zerowrapper.NewScope("websocket-server"),
	}
	w.Conn.SetPongHandler(func(appData string) error {
		_ = w.Conn.SetReadDeadline(time.Now().Add(w.PongWait))
		return nil
	})
	return w, nil
}

func (w *WsServer) Start() {
	defer func() {
		w.logger.Info().Msg("disconnected")
		_ = w.Conn.Close()
	}()
	for {
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			w.logger.Error().Msgf("read:%v", err)
			return
		}
		w.logger.Info().Msgf("recv req: %s", message)
		req := protoo.Message{}
		d := json.NewDecoder(bytes.NewReader(message))
		d.UseNumber()
		if err := d.Decode(&req); err != nil {
			w.logger.Error().Msgf("unmarshal request error%v", err)
			continue
		}
		var rsp *protoo.Message
		if req.Request {
			rsp = w.RequestHandler(req)
		}

		w.logger.Info().Msgf("send rsp: %s", rsp)
		err = w.Conn.WriteMessage(mt, rsp.Marshal())
		if err != nil {
			w.logger.Error().Msgf("write:%v", err)
			continue
		}
	}
}
