package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/byyam/mediasoup-go-worker/example/server/webrtctransport"
	"github.com/byyam/mediasoup-go-worker/global"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/gorilla/websocket"
	"github.com/jiyeyuran/go-protoo"
)

var (
	logger = utils.NewLogger("server")
	worker *mediasoup_go_worker.SimpleWorker
)

func main() {
	go func() {
		http.HandleFunc("/echo", echo)
		log.Fatal(http.ListenAndServe("localhost:8080", nil))
	}()
	conf.InitCli()
	worker = mediasoup_go_worker.NewSimpleWorker()
	worker.Start()
	// block here
	listenSignal()
	worker.Stop()
}

func echo(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // use default options
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:%v", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	h := webrtctransport.NewHandler(worker)
	s := wsconn.NewWsServer(wsconn.WsServerOpt{
		PingInterval: 10 * time.Second,
		PongWait:     1 * time.Minute,
		Conn:         c,
		Handlers: map[string]func(protoo.Message) *protoo.Message{
			isignal.MethodPublish:   h.PublishHandler,
			isignal.MethodUnPublish: h.UnPublishHandler,
		},
	})
	s.Start()
}

func listenSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals
	logger.Warn("[pid=%d]stop worker", global.Pid)
}
