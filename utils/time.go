package utils

import (
	"time"
)

const (
	millisecond = int64(time.Millisecond)
)

// GetMilliseconds get timestamp with milliseconds
func GetMilliseconds(t ...time.Time) int64 {
	if len(t) < 1 {
		return time.Now().UnixNano() / millisecond
	}

	return t[0].UnixNano() / millisecond
}
