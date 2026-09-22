package algo

import (
	"cmp"
	"slices"
)

// Sort sorts a in ascending order (in place).
func Sort[T cmp.Ordered](a []T) {
	slices.Sort(a)
}

// SortFunc sorts a using cmpFn (in place). cmpFn(a, b) should return <0 when a < b.
func SortFunc[T any](a []T, cmpFn func(a, b T) int) {
	slices.SortFunc(a, cmpFn)
}

// SortStable stably sorts a using cmpFn (in place).
func SortStable[T any](a []T, cmpFn func(a, b T) int) {
	slices.SortStableFunc(a, cmpFn)
}

// IsSorted reports whether a is sorted in ascending order.
func IsSorted[T cmp.Ordered](a []T) bool {
	return slices.IsSorted(a)
}
