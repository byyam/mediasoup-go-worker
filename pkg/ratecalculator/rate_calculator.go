package ratecalculator

import (
	"math"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
)

const (
	DefaultWindowSize  = 1000
	DefaultBpsScale    = 8000
	DefaultWindowItems = 100
)

type BufferItem struct {
	count int
	time  int64
}

func (p *BufferItem) Reset() {
	p.count = 0
	p.time = 0
}

type RateCalculator struct {
	windowSizeMs        int          // Window Size (in milliseconds).
	scale               float64      // Scale in which the rate is represented.
	windowItems         int32        // Window Size (number of items).
	itemSizeMs          int64        // Item Size (in milliseconds), calculated as: windowSizeMs / windowItems.
	buffer              []BufferItem // Buffer to keep data.
	newestItemStartTime int64        // Time (in milliseconds) for last item in the time window.
	newestItemIndex     int32        // Index for the last item in the time window.
	oldestItemStartTime int64        // Time (in milliseconds) for oldest item in the time window.
	oldestItemIndex     int32        // Index for the oldest item in the time window.
	totalCount          int          // Total count in the time window.
	bytes               int64        // Total bytes transmitted.
	lastRate            uint32       // Last value calculated by GetRate().
	lastTime            int64        // Last time GetRate() was called.

	logger utils.Logger
}

func (p RateCalculator) GetBytes() int64 {
	return p.bytes
}

func (p *RateCalculator) GetRate(nowMs int64) uint32 {
	if nowMs == p.lastTime {
		return p.lastRate
	}
	p.RemoveOldData(nowMs)

	scale := p.scale / float64(p.windowSizeMs)
	p.lastTime = nowMs
	p.lastRate = uint32(math.Trunc(float64(p.totalCount)*scale + 0.5))

	return p.lastRate
}

func (p *RateCalculator) Update(size int, nowMs int64) {
	// Ignore too old data. Should never happen.
	if nowMs < p.oldestItemStartTime {
		return
	}
	// Increase bytes.
	p.bytes += int64(size)
	p.RemoveOldData(nowMs)

	// If the elapsed time from the newest item start time is greater than the
	// item size (in milliseconds), increase the item index.
	if p.newestItemIndex < 0 || nowMs-p.newestItemStartTime >= p.itemSizeMs {
		p.newestItemIndex++
		p.newestItemStartTime = nowMs
		if p.newestItemIndex >= p.windowItems {
			p.newestItemIndex = 0
		}
		// Newest index overlaps with the oldest one, remove it.
		if p.newestItemIndex == p.oldestItemIndex && p.oldestItemIndex != -1 {
			p.logger.Warn("calculation buffer full, windowSizeMs:%d ms windowItems:%d", p.windowSizeMs, p.windowItems)

			oldestItem := &p.buffer[p.oldestItemIndex]
			p.totalCount -= oldestItem.count
			oldestItem.Reset()
			if p.oldestItemIndex+1 >= p.windowItems {
				p.oldestItemIndex = 0
			} else {
				p.oldestItemStartTime += 1
			}
		}
		// Set the newest item.
		item := &p.buffer[p.newestItemIndex]
		item.count = size
		item.time = nowMs
	} else {
		// Update the newest item.
		item := &p.buffer[p.newestItemIndex]
		item.count += size
	}
	// Set the oldest item index and time, if not set.
	if p.oldestItemIndex < 0 {
		p.oldestItemIndex = p.newestItemIndex
		p.oldestItemStartTime = nowMs
	}
	p.totalCount += size
	// Reset lastRate and lastTime so GetRate() will calculate rate again even
	// if called with same now in the same loop iteration.
	p.lastRate = 0
	p.lastTime = 0
}

func (p *RateCalculator) RemoveOldData(nowMs int64) {
	p.logger.Trace("remove old data%+v,nowMs=%d", *p, nowMs)

	// No item set.
	if p.newestItemIndex < 0 || p.oldestItemIndex < 0 {
		return
	}
	newOldestTime := nowMs - int64(p.windowSizeMs)
	// Oldest item already removed.
	if newOldestTime <= p.oldestItemStartTime {
		p.logger.Trace("oldest item already removed")
		return
	}
	// A whole window size time has elapsed since last entry. Reset the buffer.
	if newOldestTime > p.newestItemStartTime {
		p.Reset()
		return
	}

	for p.oldestItemStartTime < newOldestTime {
		p.logger.Trace("oldestItemStartTime=%d,newOldestTime=%d", p.oldestItemStartTime, newOldestTime)
		oldestItem := p.buffer[p.oldestItemIndex]
		p.totalCount -= oldestItem.count
		oldestItem.Reset()

		if p.oldestItemIndex+1 >= p.windowItems {
			p.oldestItemIndex = 0
		} else {
			p.oldestItemIndex += 1
		}

		newOldestItem := p.buffer[p.oldestItemIndex]
		p.oldestItemStartTime = newOldestItem.time
		p.logger.Trace("update oldestItemStartTime=%d", p.oldestItemStartTime)
	}
}

func (p *RateCalculator) Reset() {
	p.logger.Trace("reset %+v", *p)
	p.newestItemStartTime = 0
	p.newestItemIndex = -1
	p.oldestItemStartTime = 0
	p.oldestItemIndex = -1
	p.totalCount = 0
	p.lastRate = 0
	p.lastTime = 0
	for _, b := range p.buffer {
		b.Reset()
	}
}

func NewRateCalculator(windowSizeMs int, scale float64, windowItems int32) *RateCalculator {
	if windowSizeMs == 0 {
		windowSizeMs = DefaultWindowSize
	}
	if scale == 0 {
		scale = DefaultBpsScale
	}
	if windowItems == 0 {
		windowItems = DefaultWindowItems
	}

	return &RateCalculator{
		windowSizeMs:    windowSizeMs,
		scale:           scale,
		windowItems:     windowItems,
		itemSizeMs:      int64(math.Max(float64(windowSizeMs)/float64(windowItems), 1)),
		buffer:          make([]BufferItem, windowItems),
		newestItemIndex: -1,
		oldestItemIndex: -1,
		logger:          utils.NewLogger("rate-calculator"),
	}
}
