package main

import (
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/pion/rtp"
	"log"
	"net"
)

const (
	protocol = "udp"
	addr     = "10.12.165.74:50002"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr(protocol, addr)
	checkErr(err)

	udpSocket, err := net.ListenUDP(protocol, udpAddr)
	checkErr(err)
	defer func() {
		_ = udpSocket.Close()
	}()

	buf := make([]byte, global.ReceiveMTU)
	rtpPacket := &rtp.Packet{}
	for {
		n, addr, err := udpSocket.ReadFromUDPAddrPort(buf)
		if err != nil {
			log.Printf("udpSocketPacketReceived error:%s", err.Error())
			continue
		}
		if err = rtpPacket.Unmarshal(buf[:n]); err != nil {
			log.Printf("rtpPacket.Unmarshal error:%v", err)
			continue
		}
		log.Printf("udpSocketPacketReceived addr:[%s][%d] %s", addr.String(), n, rtpPacket.String())
	}
}
