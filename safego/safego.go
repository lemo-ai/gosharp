// Package safego runs goroutines with panic recovery.
package safego

import (
	"fmt"
	"runtime/debug"
)

// PanicHandler is called when a recovered panic is observed.
type PanicHandler func(recovered any, stack []byte)

var defaultHandler PanicHandler = func(recovered any, stack []byte) {
	fmt.Printf("safego: panic recovered: %v\n%s\n", recovered, stack)
}

// SetPanicHandler overrides the default panic logger.
func SetPanicHandler(h PanicHandler) {
	if h != nil {
		defaultHandler = h
	}
}

// Go starts fn in a new goroutine and recovers panics via the default handler.
func Go(fn func()) {
	GoWithHandler(fn, defaultHandler)
}

// GoWithHandler starts fn in a new goroutine and recovers panics via h.
func GoWithHandler(fn func(), h PanicHandler) {
	if fn == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if h != nil {
					h(r, debug.Stack())
				}
			}
		}()
		fn()
	}()
}
