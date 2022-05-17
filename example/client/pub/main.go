package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/byyam/mediasoup-go-worker/example/internal/sdk"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/srtp/v2"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/h264reader"
	"github.com/urfave/cli/v2"
)

var (
	wsAddr = "localhost:8080"
	wsPath = "/access/jwt"

	peerId = "pub-peer"
	roomId = "demoRoom"
	appId  = "3"

	videoFileName     = "/Users/yourname/Downloads/countdown.h264"
	streamLoop        = 1
	h264FrameDuration = time.Millisecond * 33
	rtpOutboundMTU    = uint16(1200)
	clockRate         = uint32(90000)
	ssrc              = uint32(567567)
	payloadType       = uint8(125)
	streamId          = uint64(123123)
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
			&cli.UintFlag{
				Name:        "rtpOutboundMTU",
				Usage:       "Use a rtpOutboundMTU",
				Value:       uint(rtpOutboundMTU),
				DefaultText: strconv.Itoa(int(rtpOutboundMTU)),
			},
			&cli.UintFlag{
				Name:        "streamLoop",
				Usage:       "Use a streamLoop",
				Value:       uint(streamLoop),
				DefaultText: strconv.Itoa(streamLoop),
			},
			&cli.StringFlag{Name: "peerId", Required: true},
			&cli.StringFlag{Name: "roomId", Required: true},
			&cli.StringFlag{Name: "appId", Required: true},
			&cli.StringFlag{Name: "videoFileName", Required: true},
			&cli.UintFlag{Name: "h264FrameDuration", Required: true},
			&cli.UintFlag{Name: "clockRate", Required: true},
			&cli.UintFlag{Name: "ssrc", Required: true},
			&cli.UintFlag{Name: "payloadType", Required: true},
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
		if c.Uint("rtpOutboundMTU") != 0 {
			rtpOutboundMTU = uint16(c.Uint("rtpOutboundMTU"))
		}
		if c.Uint("streamLoop") != 0 {
			streamLoop = int(c.Uint("streamLoop"))
		}
		peerId = c.String("peerId")
		roomId = c.String("roomId")
		appId = c.String("appId")
		videoFileName = c.String("videoFileName")
		h264FrameDuration = time.Duration(c.Uint("h264FrameDuration")) * time.Millisecond
		clockRate = uint32(c.Uint("clockRate"))
		ssrc = uint32(c.Uint("ssrc"))
		payloadType = uint8(c.Uint("payloadType"))
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

	if _, err := os.Stat(videoFileName); os.IsNotExist(err) {
		panic(err)
	}

	// connect to sfu
	opt := sdk.PublishOpt{
		StreamId:    streamId,
		SSRC:        ssrc,
		PayloadType: payloadType,
		MimeType:    webrtc.MimeTypeH264,
		ClockRate:   clockRate,
	}
	newClient := sdk.NewClient(wsconn.WsClientOpt{
		Addr: wsAddr,
		Path: wsPath,
	})
	clientConn, err := newClient.Publish(opt)
	defer func() {
		_ = clientConn.Close()
	}()
	if err != nil {
		panic(err)
	}
	log.Println("sfu connected...")

	srtpConfig, err := newClient.GetSRTPConfig()
	if err != nil {
		panic(err)
	}
	srtpSession, err := srtp.NewSessionSRTP(clientConn, srtpConfig)
	if err != nil {
		panic(err)
	}
	rtpWriteStream, err := srtpSession.OpenWriteStream()
	if err != nil {
		panic(err)
	}
	log.Println("srtp session ready...")

	packetizer := rtp.NewPacketizer(rtpOutboundMTU, payloadType, ssrc, &codecs.H264Payloader{}, rtp.NewRandomSequencer(), clockRate)

	for i := 0; i < streamLoop; i++ {
		sendRtp(i, rtpWriteStream, packetizer)
		log.Println("sendRtp loop:", i)
	}

	_ = newClient.UnPublish(streamId)
	log.Println("finished")
}

func sendRtp(i int, rtpWriteStream *srtp.WriteStreamSRTP, packetizer rtp.Packetizer) {
	// Open a H264 file and start reading using our IVFReader
	file, h264Err := os.Open(videoFileName)
	if h264Err != nil {
		panic(h264Err)
	}
	defer func() {
		_ = file.Close()
	}()

	h264, h264Err := h264reader.NewReader(file)
	if h264Err != nil {
		panic(h264Err)
	}

	ticker := time.NewTicker(h264FrameDuration)
	for ; true; <-ticker.C {
		nal, h264Err := h264.NextNAL()
		if h264Err == io.EOF {
			log.Printf("All video frames parsed and sent, exited")
			break
		}
		if h264Err != nil {
			panic(h264Err)
		}

		packages := packetizer.Packetize(nal.Data, uint32(h264FrameDuration.Seconds()*float64(clockRate)))
		for _, rtpPkg := range packages {
			_, err := rtpWriteStream.WriteRTP(&rtpPkg.Header, rtpPkg.Payload)
			if err != nil {
				log.Println("write rtp error:", err)
				continue
			}
			// log.Printf("[%d]rtp out:%+v,len:%d", i, rtpPkg.Header, len(rtpPkg.Payload))
		}
	}
	ticker.Stop()
}
