package main

import (
	"github.com/pion/rtp"
	"log"
	"net"
	"time"
)

const (
	protocol   = "udp"
	addr       = "127.0.0.1:40001"
	remoteAddr = "127.0.0.1:40000"
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

	remoteUdpAddr, err := net.ResolveUDPAddr(protocol, remoteAddr)
	checkErr(err)

	rtpPacket := &rtp.Packet{}

	ticker := time.NewTicker(time.Second * 1)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			buf, err := rtpPacket.Marshal()
			if err != nil {
				log.Printf("rtpPacket.Marshal error:%v", err)
				continue
			}
			_, err = udpSocket.WriteToUDP(buf, remoteUdpAddr)
			if err != nil {
				log.Printf("udpSocket.WriteToUDP error:%v", err)
				continue
			}
		}
	}
}
