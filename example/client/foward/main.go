package main

import (
	"github.com/byyam/mediasoup-go-worker/example/internal/forwardpkg"
	"io"
	"log"
	"os"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"

	"github.com/pion/webrtc/v3/pkg/media/h264reader"
)

// ffplay -i h264.sdp -protocol_whitelist file,udp,rtp,rtcp

const (
	videoFileName     = "/Users/yourname/Downloads/countdown.h264"
	h264FrameDuration = time.Millisecond * 33
	rtpOutboundMTU    = 1200
	clockRate         = 90000
	ssrc              = 567567
	payloadType       = 96
	local             = "127.0.0.1:"
	remote            = "127.0.0.1:12002"
)

func main() {
	forwardInfo := forwardpkg.ForwardInfo{}
	if err := forwardInfo.InitConn(local, remote); err != nil {
		panic(err)
	}
	defer func() {
		_ = forwardInfo.Close()
	}()

	if _, err := os.Stat(videoFileName); os.IsNotExist(err) {
		panic(err)
	}
	// Open a H264 file and start reading using our IVFReader
	file, h264Err := os.Open(videoFileName)
	if h264Err != nil {
		panic(h264Err)
	}

	h264, h264Err := h264reader.NewReader(file)
	if h264Err != nil {
		panic(h264Err)
	}

	packetizer := rtp.NewPacketizer(rtpOutboundMTU, payloadType, ssrc, &codecs.H264Payloader{}, rtp.NewRandomSequencer(), clockRate)
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

		packages := packetizer.Packetize(nal.Data, uint32(h264FrameDuration.Seconds()*clockRate))
		for _, rtpPkg := range packages {
			forwardInfo.Forward(rtpPkg)
		}
	}
	log.Println("finished")
}
