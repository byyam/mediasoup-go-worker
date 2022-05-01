package seqmgr

import "math"

func IsSeqHigherThanUint32(lhs, rhs uint32) bool {
	maxValue := math.MaxUint32
	return isSeqHigherThan(uint64(lhs), uint64(rhs), uint64(maxValue))
}

func IsSeqLowerThanUint32(lhs, rhs uint32) bool {
	maxValue := math.MaxUint32
	return isSeqLowerThan(uint64(lhs), uint64(rhs), uint64(maxValue))
}

func isSeqHigherThan(lhs, rhs, maxValue uint64) bool {
	if ((lhs > rhs) && lhs-rhs <= maxValue/2) || ((rhs > lhs) && rhs-lhs > maxValue/2) {
		return true
	}
	return false
}

func isSeqLowerThan(lhs, rhs, maxValue uint64) bool {
	if ((rhs > lhs) && rhs-lhs <= maxValue/2) || ((lhs > rhs) && lhs-rhs > maxValue/2) {
		return true
	}
	return false
}
