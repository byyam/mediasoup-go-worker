package nack

import (
	"time"

	"github.com/pion/rtcp"
)

const (
	maxTries      = 5                      // Max number of times a packet will be NACKed
	cacheSize     = 100                    // Max NACK sn the sfu will keep reference
	minInterval   = 20 * time.Millisecond  // minimum interval between NACK tries for the same sequence number
	maxInterval   = 400 * time.Millisecond // maximum interval between NACK tries for the same sequence number
	initialDelay  = 10 * time.Millisecond  // delay before NACKing a sequence number to allow for out-of-order packets
	backoffFactor = float64(1.25)
)

type NackQueue struct {
	nacks []*nack
	rtt   uint32
}

func NewNACKQueue() *NackQueue {
	return &NackQueue{
		nacks: make([]*nack, 0, cacheSize),
	}
}

func (n *NackQueue) SetRTT(rtt uint32) {
	n.rtt = rtt
}

func (n *NackQueue) Remove(sn uint16) {
	for idx, nack := range n.nacks {
		if nack.seqNum != sn {
			continue
		}

		copy(n.nacks[idx:], n.nacks[idx+1:])
		n.nacks = n.nacks[:len(n.nacks)-1]
		break
	}
}

func (n *NackQueue) Push(sn uint16) {
	// if at capacity, pop the first one
	if len(n.nacks) == cap(n.nacks) {
		copy(n.nacks[0:], n.nacks[1:])
		n.nacks = n.nacks[:len(n.nacks)-1]
	}

	n.nacks = append(n.nacks, newNack(sn))
}

func (n *NackQueue) Pairs() ([]rtcp.NackPair, int) {
	if len(n.nacks) == 0 {
		return nil, 0
	}

	now := time.Now()

	// set it far back to get the first pair
	baseSN := n.nacks[0].seqNum - 17

	snsToPurge := make([]uint16, 0)

	numSeqNumsNacked := 0
	isPairActive := false
	var np rtcp.NackPair
	var nps []rtcp.NackPair
	for _, nack := range n.nacks {
		shouldSend, shouldRemove, sn := nack.getNack(now, n.rtt)
		if shouldRemove {
			snsToPurge = append(snsToPurge, sn)
			continue
		}
		if !shouldSend {
			continue
		}

		numSeqNumsNacked++
		if (sn - baseSN) > 16 {
			// need a new nack pair
			if isPairActive {
				nps = append(nps, np)
				isPairActive = false
			}

			baseSN = sn

			np.PacketID = sn
			np.LostPackets = 0

			isPairActive = true
		} else {
			np.LostPackets |= 1 << (sn - baseSN - 1)
		}
	}

	// add any left over
	if isPairActive {
		nps = append(nps, np)
	}

	for _, sn := range snsToPurge {
		n.Remove(sn)
	}

	return nps, numSeqNumsNacked
}
