package util

import (
	"time"
)

// MillisecondsToTime converts milliseconds (int64) to time.Time
func MillisecondsToTime(ms int64) time.Time {
	if ms == 0 {
		return time.Time{}
	}
	return time.Unix(0, ms*int64(time.Millisecond))
}

// TimeToMilliseconds converts time.Time to milliseconds (int64)
func TimeToMilliseconds(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / int64(time.Millisecond)
}