package main

import (
	"github.com/byyam/mediasoup-go-worker/example/internal/demoutils"
	"github.com/byyam/mediasoup-go-worker/example/server/pipetransport"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byyam/mediasoup-go-worker/example/server/workerapi"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/byyam/mediasoup-go-worker/example/server/webrtctransport"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/gorilla/websocket"
)

var (
	logger = utils.NewLogger("server")
	worker *mediasoup_go_worker.SimpleWorker
)

const (
	localWsAddr                       = ":12001"
	localHttpAddr                     = ":12002"
	pathWebrtcTransport               = "/webrtc_transport"
	pathPipeTransportCreateAndConnect = "/pipe_transport/create_and_connect"
)

func main() {
	conf.InitCli()
	logger.Info("argv:%+v", conf.Settings)
	if conf.Settings.PrometheusPort > 0 {
		monitor.InitPrometheus(monitor.WithPath(conf.Settings.PrometheusPath), monitor.WithPort(conf.Settings.PrometheusPort))
	}

	worker = mediasoup_go_worker.NewSimpleWorker()
	worker.Start()
	if err := workerapi.CreateRouter(worker, demoutils.GetRouterId(worker)); err != nil {
		panic(err)
	}

	go func() {
		http.HandleFunc(pathWebrtcTransport, handleWebrtcTransport)
		log.Fatal(http.ListenAndServe(localWsAddr, nil))
	}()

	go func() {
		h := pipetransport.NewHandler(worker)
		server := http.Server{
			Addr:           localHttpAddr,
			Handler:        &myHandler{},
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		mux = make(map[string]func(http.ResponseWriter, *http.Request))
		mux[pathPipeTransportCreateAndConnect] = h.HandlePipeTransportCreateAndConnect
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// block here
	listenSignal()
	worker.Stop()
}

func handleWebrtcTransport(w http.ResponseWriter, r *http.Request) {
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
	s, err := wsconn.NewWsServer(wsconn.WsServerOpt{
		PingInterval:   10 * time.Second,
		PongWait:       1 * time.Minute,
		Conn:           c,
		RequestHandler: h.HandleProtooMessage,
	})
	if err != nil {
		panic(err)
	}
	s.Start()
}

func listenSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals
	logger.Warn("[pid=%d]stop worker", global.Pid)
}
