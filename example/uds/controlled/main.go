package main

import (
	"fmt"
	"os"

	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
	"github.com/byyam/mediasoup-go-worker/pkg/udschannel"
)

func main() {
	_, err := netparser.NewNetStringsFd(udschannel.CustomerPipeStart, udschannel.CustomerPipeStart+1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("controlled pid:%d start", os.Getpid())
	select {}
}
