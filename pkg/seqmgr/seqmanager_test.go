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
