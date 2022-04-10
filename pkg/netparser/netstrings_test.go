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
	buffer := make([]byte, 4194308)
	buf := bytes.NewReader(netstringPayload)
	parser := NewNetStrings(nil, buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ReadBuffer(buffer)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset(netstringPayload)
	}
}

func TestNetStrings_ReadBuffer(t *testing.T) {
	buffer := make([]byte, 4194308)
	buf := bytes.NewReader(netstringPayload)
	parser := NewNetStrings(nil, buf)
	n, err := parser.ReadBuffer(buffer)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buffer[:n])
	buf.Reset(netstringPayload)
	n1, err := parser.ReadBuffer(buffer)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buffer[:n1])
}
