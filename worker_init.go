package mediasoup_go_worker

import (
	"fmt"

	"github.com/hashicorp/go-version"

	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger = zerowrapper.NewScope("mediasoup-worker-init")
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func InitWorker(mediasoupVersion string) (*workerchannel.Channel, *workerchannel.PayloadChannel, error) {
	var err error
	//defer func() {
	//	if r := recover(); r != nil {
	//		logger.Error("init worker panic:%s", debug.Stack())
	//	}
	//}()

	currentLatest, err := version.NewVersion(mediasoupVersion)
	checkError(err)
	logger.Info().Msgf("MEDIASOUP_VERSION:%s", mediasoupVersion)

	// prepare write/read channel
	var netParser netparser.INetParser
	nativeJsonVersion, _ := version.NewVersion(workerchannel.NativeJsonVersion)
	nativeVersion, _ := version.NewVersion(workerchannel.NativeVersion)
	jsonFormat := true
	if currentLatest.GreaterThanOrEqual(nativeJsonVersion) {
		order := netparser.HostByteOrder()
		netParser, err = netparser.NewNetNativeFd(workerchannel.ProducerChannelFd, workerchannel.ConsumerChannelFd, order)
		logger.Info().Msgf("create native codec, host order:%s", order)
		// https://github.com/versatica/mediasoup/pull/870
		if currentLatest.GreaterThanOrEqual(nativeVersion) {
			jsonFormat = false
		}
	} else {
		netParser, err = netparser.NewNetStringsFd(workerchannel.ProducerChannelFd, workerchannel.ConsumerChannelFd)
		logger.Info().Msg("create netstrings codec")
	}
	checkError(err)

	defer func() {
		_ = netParser.Close()
	}()

	channel := workerchannel.NewChannel(netParser, fmt.Sprintf("cfd=%d,pfd=%d", workerchannel.ConsumerChannelFd, workerchannel.ProducerChannelFd), jsonFormat)
	payloadChannel := workerchannel.NewPayloadChannel()
	return channel, payloadChannel, nil
}
