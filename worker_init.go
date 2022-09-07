package mediasoup_go_worker

import (
	"fmt"

	"github.com/hashicorp/go-version"

	"github.com/byyam/mediasoup-go-worker/internal/constant"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/pkg/netparser"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

var (
	logger = utils.NewLogger("mediasoup-worker-init")
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
	logger.Info("MEDIASOUP_VERSION:%s", mediasoupVersion)

	// prepare write/read channel
	var netParser netparser.INetParser
	nativeJsonVersion, _ := version.NewVersion(constant.NativeJsonVersion)
	nativeVersion, _ := version.NewVersion(constant.NativeVersion)
	jsonFormat := true
	if currentLatest.GreaterThanOrEqual(nativeJsonVersion) {
		order := netparser.HostByteOrder()
		netParser, err = netparser.NewNetNativeFd(constant.ProducerChannelFd, constant.ConsumerChannelFd, order)
		logger.Info("create native codec, host order:%s", order)
		// https://github.com/versatica/mediasoup/pull/870
		if currentLatest.GreaterThanOrEqual(nativeVersion) {
			jsonFormat = false
		}
	} else {
		netParser, err = netparser.NewNetStringsFd(constant.ProducerChannelFd, constant.ConsumerChannelFd)
		logger.Info("create netstrings codec")
	}
	checkError(err)

	defer func() {
		_ = netParser.Close()
	}()

	channel := workerchannel.NewChannel(netParser, fmt.Sprintf("pid=%d,cfd=%d,pfd=%d", global.Pid, constant.ConsumerChannelFd, constant.ProducerChannelFd), jsonFormat)
	payloadChannel := workerchannel.NewPayloadChannel()
	return channel, payloadChannel, nil
}
