// Package collection provides generic slice helpers.
package collection

// Contains reports whether item is present in s.
func Contains[T comparable](s []T, item T) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

// Index returns the first index of item in s, or -1.
func Index[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

// Unique returns a new slice with duplicates removed, preserving first-seen order.
func Unique[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	seen := make(map[T]struct{}, len(s))
	out := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// Filter returns elements for which keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
	if len(s) == 0 {
		return nil
	}
	out := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Map applies fn to every element and returns the results.
func Map[T any, R any](s []T, fn func(T) R) []R {
	if len(s) == 0 {
		return nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

// Chunk splits s into chunks of at most size. size must be > 0.
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 || len(s) == 0 {
		return nil
	}
	n := (len(s) + size - 1) / size
	out := make([][]T, 0, n)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
	}
	return out
}

// Reverse returns a reversed copy of s.
func Reverse[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	out := make([]T, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

// Diff returns elements in a that are not in b.
func Diff[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return nil
	}
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	out := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := set[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}

// Intersect returns elements present in both a and b (order follows a).
func Intersect[T comparable](a, b []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	out := make([]T, 0)
	seen := make(map[T]struct{})
	for _, v := range a {
		if _, ok := set[v]; !ok {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
