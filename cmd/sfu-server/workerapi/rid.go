package workerapi

import (
	"sync/atomic"
)

var rid uint32

func GetRid() uint32 {
	return atomic.AddUint32(&rid, 1)
}
