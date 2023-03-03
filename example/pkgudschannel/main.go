package main

import (
	"fmt"

	"github.com/byyam/mediasoup-go-worker/pkg/udschannel"
)

func OptsMarshal() []string {
	var options = make([]string, 0, 20)
	options = append(options, fmt.Sprintf("--logLevel=%s", "wrn"))
	options = append(options, fmt.Sprintf("--logTags=%s", "rtc"))
	options = append(options, fmt.Sprintf("--rtcMinPort=%d", 40000))
	options = append(options, fmt.Sprintf("--rtcMaxPort=%d", 41000))
	options = append(options, fmt.Sprintf("--rtcStaticPort=%d", 50000))
	options = append(options, fmt.Sprintf("--rtcListenIp=%s", "127.0.0.1"))

	return options
}

func main() {
	opts := OptsMarshal()
	_, err := udschannel.NewUdsChannel("./bin/mediasoup-worker-darwin", opts, 2)
	if err != nil {
		panic(err)
	}
	select {}
}
