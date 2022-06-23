package nack

import (
	"math"
	"time"
)

type nack struct {
	seqNum       uint16
	tries        uint8
	lastNackedAt time.Time
}

func newNack(sn uint16) *nack {
	return &nack{
		seqNum:       sn,
		tries:        0,
		lastNackedAt: time.Now(),
	}
}

func (n *nack) getNack(now time.Time, rtt uint32) (shouldSend bool, shouldRemove bool, sn uint16) {
	sn = n.seqNum
	if n.tries >= maxTries {
		shouldRemove = true
		return
	}

	var requiredInterval time.Duration
	if n.tries > 0 {
		// exponentially backoff retries, but cap maximum spacing between retries
		requiredInterval = maxInterval
		backoffInterval := time.Duration(float64(rtt)*math.Pow(backoffFactor, float64(n.tries-1))) * time.Millisecond
		if backoffInterval < requiredInterval {
			requiredInterval = backoffInterval
		}
	}
	if requiredInterval < minInterval {
		//
		// Wait for some time for out-of-order packets before NACKing even if before NACKing first time.
		// For subsequent tries, maintain minimum spacing.
		//
		requiredInterval = minInterval
	}

	if now.Sub(n.lastNackedAt) < requiredInterval {
		return
	}

	n.tries++
	n.lastNackedAt = now
	shouldSend = true
	return
}
