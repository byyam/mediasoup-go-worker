package wsconn

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/jiyeyuran/go-protoo"

	"github.com/gorilla/websocket"

	"github.com/byyam/mediasoup-go-worker/utils"
)

type WsServer struct {
	logger utils.Logger
	WsServerOpt
}

type WsServerOpt struct {
	PingInterval time.Duration
	PongWait     time.Duration
	Conn         *websocket.Conn
	Handlers     map[string]func(protoo.Message) *protoo.Message
}

func NewWsServer(opt WsServerOpt) *WsServer {
	w := &WsServer{
		WsServerOpt: opt,
		logger:      utils.NewLogger("websocket-server"),
	}

	w.Conn.SetPongHandler(func(appData string) error {
		_ = w.Conn.SetReadDeadline(time.Now().Add(w.PongWait))
		return nil
	})

	return w
}

func (w *WsServer) Start() {
	defer func() {
		w.logger.Info("disconnected")
		_ = w.Conn.Close()
	}()
	for {
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			w.logger.Error("read:%v", err)
			return
		}
		w.logger.Info("recv: %s", message)
		req := protoo.Message{}
		d := json.NewDecoder(bytes.NewReader(message))
		d.UseNumber()
		if err := d.Decode(&req); err != nil {
			w.logger.Error("unmarshal request error%v", err)
			continue
		}
		fn := w.Handlers[req.Method]
		rsp := fn(req)
		err = w.Conn.WriteMessage(mt, rsp.Marshal())
		if err != nil {
			w.logger.Error("write:%v", err)
			continue
		}
	}
}
