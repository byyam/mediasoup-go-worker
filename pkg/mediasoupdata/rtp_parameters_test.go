package mediasoupdata

import (
	"fmt"
	"testing"

	FBS__RtpParameters "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpParameters"
)

func Test_RtpCodecSpecificParameters_Set(t *testing.T) {
	params := RtpCodecSpecificParameters{}
	fbs := make([]*FBS__RtpParameters.ParameterT, 0)
	aptValue := map[string]int32{
		"value": 96,
	}
	apt := &FBS__RtpParameters.ParameterT{
		Name: "apt",
		Value: &FBS__RtpParameters.ValueT{
			Type:  FBS__RtpParameters.ValueInteger32,
			Value: aptValue,
		},
	}
	fbs = append(fbs, apt)
	params.Set(fbs)
	fmt.Printf("%+v\n", params)
}

func Test_RtpCodecSpecificParameters_Convert(t *testing.T) {
	params := &RtpCodecSpecificParameters{
		Apt: 96,
	}
	fbs := params.Convert()

	fmt.Printf("%+v\n", fbs)
}
