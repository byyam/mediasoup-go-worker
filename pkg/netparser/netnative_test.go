package netparser

import (
	"encoding/binary"
	"io"
	"testing"
)

func BenchmarkNetNative_WriteBuffer(b *testing.B) {
	parser := NewNetNative(io.Discard, nil, HostByteOrder())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.WriteBuffer(rawPayload)
	}
}

func BenchmarkNetNative_ReadBuffer(b *testing.B) {
	r, w := io.Pipe()

	parser := NewNetNative(nil, r, HostByteOrder())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resetWriter(w, rawPayload)
		_, err := parser.ReadBuffer()
		if err != nil {
			b.Fatal(err)
		}
		// b.Logf("out[%d]:%s", len(out), string(out))
	}
	_ = w.Close()
	_ = r.Close()
}

func TestNetNative_ReadBuffer(t *testing.T) {
	r, w := io.Pipe()
	payload := []byte("{\"id\":2,\"internal\":null,\"method\":\"worker.getResourceUsage\"}")
	t.Logf("payload[%d]", len(payload))
	resetWriter(w, payload)

	// read
	parser := NewNetNative(nil, r, HostByteOrder())
	out, err := parser.ReadBuffer()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("out[%d]:%s", len(out), string(out))
	// reset
	resetWriter(w, payload)
	out1, err := parser.ReadBuffer()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("out1[%d]:%s", len(out1), string(out1))
	_ = w.Close()
	_ = r.Close()
}

func resetWriter(w io.Writer, payload []byte) {
	go func() {
		if err := binary.Write(w, binary.LittleEndian, uint32(len(payload))); err != nil {
			return
		}
		_, err := w.Write(payload)
		if err != nil {
			return
		}
	}()
}
