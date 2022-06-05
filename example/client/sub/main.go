package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/byyam/mediasoup-go-worker/example/internal/sdk"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/pion/rtp"
	"github.com/pion/srtp/v2"
	"github.com/urfave/cli/v2"
)

var (
	wsAddr = "localhost:8080"
	wsPath = "/echo"

	streamId = uint64(123123)
)

const (
	local  = "127.0.0.1:"
	remote = "127.0.0.1:12001"
)

func initCli() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "wsAddr",
				Usage:       "Use a local wsAddr",
				Value:       wsAddr,
				DefaultText: wsAddr,
			},
			&cli.StringFlag{
				Name:        "wsPath",
				Usage:       "Use a local wsPath",
				Value:       wsPath,
				DefaultText: wsPath,
			},
			&cli.StringFlag{Name: "peerId", Required: true},
			&cli.StringFlag{Name: "roomId", Required: true},
			&cli.StringFlag{Name: "appId", Required: true},
			&cli.Uint64Flag{Name: "streamId", Required: true},
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("wsAddr") != "" {
			wsAddr = c.String("wsAddr")
		}
		if c.String("wsPath") != "" {
			wsPath = c.String("wsPath")
		}
		//peerId = c.String("peerId")
		//roomId = c.String("roomId")
		//appId = c.String("appId")
		streamId = c.Uint64("streamId")
		return nil
	}
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		log.Fatal("Ha HA.  I am the help!!")
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initCli()

	//forwardInfo := forwardpkg.ForwardInfo{}
	//if err := forwardInfo.InitConn(local, remote); err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	_ = forwardInfo.Close()
	//}()

	opt := sdk.SubscribeOpt{
		StreamId:   streamId,
		DtlsClient: false,
	}
	newClient := sdk.NewClient(wsconn.WsClientOpt{
		Addr: wsAddr,
		Path: wsPath,
	})
	clientConn, rsp, err := newClient.Subscribe(opt)
	defer func() {
		_ = clientConn.Close()
	}()
	if err != nil {
		panic(err)
	}
	log.Println("sfu connected...", rsp)
	ssrc := rsp.RtpParameters.Encodings[0].Ssrc

	srtpConfig, err := newClient.GetSRTPConfig()
	if err != nil {
		panic(err)
	}
	srtpSession, err := srtp.NewSessionSRTP(clientConn, srtpConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = srtpSession.Close()
	}()

	go func() {
		for {
			_, ssrc, err := srtpSession.AcceptStream()
			if err != nil {
				log.Printf("accept stream error:%v", err)
			}
			log.Printf("accept new stream:%d", ssrc)
		}
	}()

	rtpReadStream, err := srtpSession.OpenReadStream(ssrc)
	if err != nil {
		panic(err)
	}
	log.Printf("srtp session ready...ssrc=%d", ssrc)

	// Receive messages in a loop from the remote peer
	// ReadConn(clientConn)
	buf := make([]byte, 1500)
	rtpPacket := &rtp.Packet{}
	go func() {
		for {
			n, header, err := rtpReadStream.ReadRTP(buf)
			if err != nil {
				panic(err)
			}
			if err = rtpPacket.Unmarshal(buf[:n]); err != nil {
				log.Printf("rtp unmarshal error:%v", err)
			}
			// forwardInfo.Forward(rtpPacket)
			log.Printf("rtp in:%+v,len=%d", header, n)
		}
	}()
	select {}
}

func ReadConn(clientConn net.Conn) { // get raw pkg without decrypted
	buf := make([]byte, 1500)
	rtpPacket := &rtp.Packet{}
	go func() {
		for {
			n, err := clientConn.Read(buf)
			if err != nil {
				panic(err)
			}
			if err = rtpPacket.Unmarshal(buf[:n]); err != nil {
				log.Printf("rtp unmarshal error:%v", err)
			}
			log.Printf("rtp[%d]:%+v", n, rtpPacket.Header)
		}
	}()
}
