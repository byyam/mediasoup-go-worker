package seqmgr

import "testing"

func TestIsSeqHigherThanUint32(t *testing.T) {
	v := IsSeqHigherThanUint32(1000, 1000)
	t.Log(v)
}

func TestIsSeqLowerThanUint32(t *testing.T) {
	v := IsSeqLowerThanUint32(1000, 1000)
	t.Log(v)
}

func TestIsSeqHigherThanUint16(t *testing.T) {
	v := IsSeqHigherThanUint16(1000, 1000)
	t.Log(v)
}

func TestIsSeqLowerThanUint16(t *testing.T) {
	v := IsSeqLowerThanUint16(1000, 1000)
	t.Log(v)
}
