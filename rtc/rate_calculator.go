package rtc

const (
	DefaultWindowSize  = 1000
	DefaultBpsScale    = 8000
	DefaultWindowItems = 100
)

type BufferItem struct {
	count int
	time  int64
}

type RateCalculator struct {
	windowSizeMs int          // Window Size (in milliseconds).
	scale        float64      // Scale in which the rate is represented.
	windowItems  uint16       // Window Size (number of items).
	itemSizeMs   int          // Item Size (in milliseconds), calculated as: windowSizeMs / windowItems.
	buffer       []BufferItem // Buffer to keep data.
	// todo
}

func newRateCalculator(windowSizeMs int, scale float64, windowItems uint16) *RateCalculator {
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
		windowSizeMs: windowSizeMs,
		scale:        scale,
		windowItems:  windowItems,
	}
}

type RtpDataCounter struct {
	rate    *RateCalculator
	packets int
}

func newRtpDataCounter(windowSizeMs int) *RtpDataCounter {
	size := 2500
	if windowSizeMs > 0 {
		size = windowSizeMs
	}
	return &RtpDataCounter{
		rate:    newRateCalculator(size, 0, 0),
		packets: 0,
	}
}
