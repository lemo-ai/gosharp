package algo

import "cmp"

// BinarySearch returns the index of target in a sorted ascending slice, or -1.
func BinarySearch[T cmp.Ordered](a []T, target T) int {
	lo, hi := 0, len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if a[mid] < target {
			lo = mid + 1
		} else if a[mid] > target {
			hi = mid
		} else {
			return mid
		}
	}
	return -1
}

// LowerBound returns the first index i in sorted a such that a[i] >= target.
// If all elements are < target, returns len(a).
func LowerBound[T cmp.Ordered](a []T, target T) int {
	lo, hi := 0, len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if a[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// UpperBound returns the first index i in sorted a such that a[i] > target.
// If all elements are <= target, returns len(a).
func UpperBound[T cmp.Ordered](a []T, target T) int {
	lo, hi := 0, len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if a[mid] <= target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
