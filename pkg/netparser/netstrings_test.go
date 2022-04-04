package netparser

import (
	"bytes"
	"io"
	"testing"
)

// go test -bench=. -benchmem -run=none

var (
	rawPayload       = []byte("{\"id\":5,\"accepted\":true,\"data\":{\"type\":\"simple\",\"rtpParameters\":{\"codecs\":null,\"rtcp\":{}},\"consumableRtpParameters\":{\"codecs\":null,\"rtcp\":{}}}}")
	netstringPayload = []byte("143:{\"id\":5,\"accepted\":true,\"data\":{\"type\":\"simple\",\"rtpParameters\":{\"codecs\":null,\"rtcp\":{}},\"consumableRtpParameters\":{\"codecs\":null,\"rtcp\":{}}}},")
)

func BenchmarkNetStrings_WriteBuffer(b *testing.B) {
	parser := NewNetStrings(io.Discard, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.WriteBuffer(rawPayload)
	}
}

func BenchmarkNetStrings_ReadBuffer(b *testing.B) {
	buf := bytes.NewReader(netstringPayload)
	parser := NewNetStrings(nil, buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ReadBuffer()
		buf.Reset(netstringPayload)
	}
}

func TestNetStrings_ReadBuffer(t *testing.T) {
	buf := bytes.NewReader(netstringPayload)
	parser := NewNetStrings(nil, buf)
	out, err := parser.ReadBuffer()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
	buf.Reset(netstringPayload)
	out1, err1 := parser.ReadBuffer()
	if err1 != nil {
		t.Fatal(err1)
	}
	t.Log(out1)
}
