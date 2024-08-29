package gerror

import "runtime"

// stack represents a stack of program counters.
type stack []uintptr

const (
	gMaxStackDepth = 32
)

// callers returns the stack callers.
func callers(skip ...int) stack {
	var (
		pcs [gMaxStackDepth]uintptr
		n   = 3
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
