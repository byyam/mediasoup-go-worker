// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WebRtcTransport

import "strconv"

type IceCandidateType byte

const (
	IceCandidateTypeHOST IceCandidateType = 0
)

var EnumNamesIceCandidateType = map[IceCandidateType]string{
	IceCandidateTypeHOST: "HOST",
}

var EnumValuesIceCandidateType = map[string]IceCandidateType{
	"HOST": IceCandidateTypeHOST,
}

func (v IceCandidateType) String() string {
	if s, ok := EnumNamesIceCandidateType[v]; ok {
		return s
	}
	return "IceCandidateType(" + strconv.FormatInt(int64(v), 10) + ")"
}
