package signaldefine

import (
	"encoding/json"
	"fmt"
	"testing"

	FBS__Producer "github.com/byyam/mediasoup-go-worker/fbs/FBS/Producer"
	FBS__RtpStream "github.com/byyam/mediasoup-go-worker/fbs/FBS/RtpStream"
)

var fbsContent = "{\"stats\":[{\"data\":{\"Type\":2,\"Value\":{\"base\":{\"data\":{\"Type\":1,\"Value\":{\"timestamp\":1724757566485,\"ssrc\":1337171516,\"kind\":0,\"mime_type\":\"audio/opus\",\"packets_lost\":0,\"fraction_lost\":0,\"packets_discarded\":0,\"packets_retransmitted\":0,\"packets_repaired\":0,\"nack_count\":0,\"nack_packet_count\":0,\"pli_count\":0,\"fir_count\":0,\"score\":10,\"rid\":\"\",\"rtx_ssrc\":0,\"rtx_packets_discarded\":0,\"round_trip_time\":0}}},\"jitter\":0,\"packet_count\":23,\"byte_count\":5676,\"bitrate\":7568,\"bitrate_by_layer\":null}}}]}"

func Test_GetProducerStatResponseSet(t *testing.T) {
	fbs := &FBS__Producer.GetStatsResponseT{
		Stats: make([]*FBS__RtpStream.StatsT, 0),
	}
	_ = json.Unmarshal([]byte(fbsContent), fbs)
	fmt.Printf("body type:%+v value:%+v\n", fbs.Stats[0].Data.Type, fbs.Stats[0].Data.Value)

	stats := GetProducerStatResponseSet(fbs)
	statsContent, _ := json.Marshal(stats)
	fmt.Printf("stats:%+v\n", string(statsContent))
}
