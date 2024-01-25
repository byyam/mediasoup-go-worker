// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package RtpStream

import "strconv"

type StatsData byte

const (
	StatsDataNONE      StatsData = 0
	StatsDataBaseStats StatsData = 1
	StatsDataRecvStats StatsData = 2
	StatsDataSendStats StatsData = 3
)

var EnumNamesStatsData = map[StatsData]string{
	StatsDataNONE:      "NONE",
	StatsDataBaseStats: "BaseStats",
	StatsDataRecvStats: "RecvStats",
	StatsDataSendStats: "SendStats",
}

var EnumValuesStatsData = map[string]StatsData{
	"NONE":      StatsDataNONE,
	"BaseStats": StatsDataBaseStats,
	"RecvStats": StatsDataRecvStats,
	"SendStats": StatsDataSendStats,
}

func (v StatsData) String() string {
	if s, ok := EnumNamesStatsData[v]; ok {
		return s
	}
	return "StatsData(" + strconv.FormatInt(int64(v), 10) + ")"
}
