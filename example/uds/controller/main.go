package main

import (
	"fmt"

	"github.com/byyam/mediasoup-go-worker/pkg/udschannel"
)

func settings() []string {
	var options = make([]string, 0)
	options = append(options, fmt.Sprintf("--logLevel=%s", "wrn"))

	return options
}

func main() {
	opts := settings()
	_, err := udschannel.NewController("../controlled/controlled", opts, 2)
	if err != nil {
		panic(err)
	}
	select {}
}
