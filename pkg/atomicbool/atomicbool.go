package atomicbool

import "sync/atomic"

type AtomicBool struct {
	val int32
}

func (b *AtomicBool) Set(value bool) { // nolint: unparam
	var i int32
	if value {
		i = 1
	}

	atomic.StoreInt32(&(b.val), i)
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&(b.val)) != 0
}

func (b *AtomicBool) Swap(value bool) bool {
	var i int32 = 0
	if value {
		i = 1
	}
	return atomic.SwapInt32(&(b.val), i) != 0
}
