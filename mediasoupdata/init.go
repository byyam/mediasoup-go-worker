package mediasoupdata

import (
	"fmt"
	"strings"
)

func Init() {
	initRtpCodecMimeType()
	initRtpCodecMimeSubType()
}

func initRtpCodecMimeType() {
	for k, v := range rtpCodecMimeType2String {
		rtpCodecMimeType2String[k] = strings.ToLower(v)
	}
	logger.Info(fmt.Sprintf("rtpCodecMimeType2String %+v", rtpCodecMimeType2String))
	for k, v := range rtpCodecMimeType2String {
		rtpCodecMimeString2Type[v] = k
	}
	logger.Info(fmt.Sprintf("rtpCodecMimeString2Type %+v", rtpCodecMimeString2Type))
}

func initRtpCodecMimeSubType() {
	for k, v := range rtpCodecMimeSubType2String {
		rtpCodecMimeSubType2String[k] = strings.ToLower(v)
	}
	logger.Info(fmt.Sprintf("rtpCodecMimeSubType2String %+v", rtpCodecMimeSubType2String))
	for k, v := range rtpCodecMimeSubType2String {
		rtpCodecMimeString2SubType[v] = k
	}
	logger.Info(fmt.Sprintf("rtpCodecMimeString2SubType %+v", rtpCodecMimeString2SubType))
}
