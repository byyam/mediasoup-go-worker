package wsconn

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jiyeyuran/go-protoo"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8080"

var upgrader = websocket.Upgrader{} // use default options

func TestNewWsServer(t *testing.T) {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	s, _ := NewWsServer(WsServerOpt{
		PingInterval: 10 * time.Second,
		PongWait:     1 * time.Minute,
		Conn:         c,
		RequestHandler: func(message protoo.Message) *protoo.Message {
			log.Printf("handle %s", isignal.MethodUnPublish)
			rspData := isignal.PublishResponse{
				TransportId: "demoId",
			}
			rsp := protoo.CreateSuccessResponse(message, rspData)
			return &rsp
		},
	})
	s.Start()
}
