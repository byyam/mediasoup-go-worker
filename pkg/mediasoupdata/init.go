package mediasoupdata

import (
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
	logger.Info().Msgf("rtpCodecMimeType2String %+v", rtpCodecMimeType2String)
	for k, v := range rtpCodecMimeType2String {
		rtpCodecMimeString2Type[v] = k
	}
	logger.Info().Msgf("rtpCodecMimeString2Type %+v", rtpCodecMimeString2Type)
}

func initRtpCodecMimeSubType() {
	for k, v := range rtpCodecMimeSubType2String {
		rtpCodecMimeSubType2String[k] = strings.ToLower(v)
	}
	logger.Info().Msgf("rtpCodecMimeSubType2String %+v", rtpCodecMimeSubType2String)
	for k, v := range rtpCodecMimeSubType2String {
		rtpCodecMimeString2SubType[v] = k
	}
	logger.Info().Msgf("rtpCodecMimeString2SubType %+v", rtpCodecMimeString2SubType)
}
