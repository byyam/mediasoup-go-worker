// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Transport

import "strconv"

type BweType byte

const (
	BweTypeTRANSPORT_CC BweType = 0
	BweTypeREMB         BweType = 1
)

var EnumNamesBweType = map[BweType]string{
	BweTypeTRANSPORT_CC: "TRANSPORT_CC",
	BweTypeREMB:         "REMB",
}

var EnumValuesBweType = map[string]BweType{
	"TRANSPORT_CC": BweTypeTRANSPORT_CC,
	"REMB":         BweTypeREMB,
}

func (v BweType) String() string {
	if s, ok := EnumNamesBweType[v]; ok {
		return s
	}
	return "BweType(" + strconv.FormatInt(int64(v), 10) + ")"
}
