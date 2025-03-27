package clock

import (
	"time"
)

func UnixMilli() int64 {
	currentTime := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	return currentTime + GlobalClockSync.GetCurrentOffset()
}

func UnixSeconds() int64 {
	return UnixMilli() / 1000
}
