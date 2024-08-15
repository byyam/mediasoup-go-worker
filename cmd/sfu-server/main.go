package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/room"
	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/monitor"
	"github.com/byyam/mediasoup-go-worker/pkg/zaplog"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/demoutils"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/pipetransport"
	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/sfuconf"

	"github.com/byyam/mediasoup-go-worker/cmd/sfu-server/workerapi"

	"github.com/gorilla/websocket"

	mediasoup_go_worker "github.com/byyam/mediasoup-go-worker"

	"github.com/byyam/mediasoup-go-worker/pkg/wsconn"
)

var (
	worker *mediasoup_go_worker.SimpleWorker

	logConfig = zerowrapper.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "sfu.log",
		MaxSize:               1,
		MaxBackups:            10,
		MaxAge:                2,
		LogTimeFieldFormat:    zerolog.TimeFormatUnixMicro,
	}
	logger zerolog.Logger
)

const (
	// wss
	localWsAddr         = ":4443"
	pathWebrtcTransport = "/"
	// tls key
	HTTPS_CERT_FULLCHAIN = "./cert/fullchain.pem"
	HTTPS_CERT_PRIVKEY   = "./cert/privkey.pem"
	// http
	localHttpAddr                     = ":12002"
	pathPipeTransportCreateAndConnect = "/pipe_transport/create_and_connect"
)

var (
	githash    string
	gitbranch  string
	buildstamp string
	goversion  string
)

func printVersion() {
	log.Printf("%11s %s", "GIT_HASH:", githash)
	log.Printf("%11s %s", "GIT_BRANCH:", gitbranch)
	log.Printf("%11s %s", "BUILD_TIME:", buildstamp)
	log.Printf("%11s %s", "GO_VERSION:", goversion)
}

// Define a map to implement routing table.
var mux map[string]func(http.ResponseWriter, *http.Request)

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement route forwarding
	if h, ok := mux[r.URL.String()]; ok {
		//Implement route forwarding with this handler, the corresponding route calls the corresponding func.
		h(w, r)
		return
	}
	_, _ = io.WriteString(w, "unknown URL: "+r.URL.String())
}

func main() {
	printVersion()
	sfuconf.InitConfig()
	zerowrapper.InitLog(logConfig)
	zaplog.Init(zaplog.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "media.log",
		MaxSize:               1,
		MaxBackups:            1,
		MaxAge:                1,
		LogTimeFieldFormat:    "",
		ErrorStackMarshaler:   false,
	})

	logger.Info().Msgf("argv:%+v", conf.Settings)
	if conf.Settings.PrometheusPort > 0 {
		monitor.InitPrometheus(monitor.WithPath(conf.Settings.PrometheusPath), monitor.WithPort(conf.Settings.PrometheusPort))
	}

	worker = mediasoup_go_worker.NewSimpleWorker()
	worker.Start()
	// default router, for pub/sub to get global router
	if err := workerapi.CreateRouter(worker, demoutils.GetRouterId(worker)); err != nil {
		panic(err)
	}

	go func() {
		http.HandleFunc(pathWebrtcTransport, handleWebrtcTransport)
		//log.Fatal(http.ListenAndServe(localWsAddr, nil))
		log.Fatal(http.ListenAndServeTLS(localWsAddr, HTTPS_CERT_FULLCHAIN, HTTPS_CERT_PRIVKEY, nil)) // tls
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
	query := r.URL.Query()
	roomId := query.Get("roomId")
	peerId := query.Get("peerId")
	if roomId == "" || peerId == "" {
		http.Error(w, "Connection request without roomId and/or peerId", http.StatusBadRequest)
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"protoo"},
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Msgf("upgrade:%v", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	h := room.NewHandler(worker, &room.QueryParams{
		RoomId: roomId,
		PeerId: peerId,
	})
	s, err := wsconn.NewWsServer(wsconn.WsServerOpt{
		TraceId:        fmt.Sprintf("%s-%s", roomId, peerId),
		PingInterval:   10 * time.Second,
		PongWait:       1 * time.Minute,
		Conn:           c,
		RequestHandler: h.HandleProtooMessage, // mediasoup proto
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
	logger.Warn().Msgf("[pid=%d]stop worker", worker.GetPid())
}
