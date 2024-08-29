package gtime

import (
	"time"
)

// wrapper is a wrapper for stdlib struct gtime.Time.
// It's used for overwriting some functions of gtime.Time, for example: String.
type wrapper struct {
	time.Time
}

// String overwrites the String function of gtime.Time.
func (t wrapper) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
