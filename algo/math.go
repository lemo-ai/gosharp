// Package algo provides common algorithm helpers.
package algo

import "cmp"

// Min returns the smaller of a and b.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp returns v limited to [lo, hi]. If lo > hi the bounds are swapped.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Abs returns the absolute value of n.
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// GCD returns the greatest common divisor of a and b (non-negative).
func GCD(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of a and b (non-negative).
func LCM(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	g := GCD(a, b)
	return Abs(a/g*b)
}

// PowInt returns base^exp for non-negative exp. Returns 1 when exp == 0.
func PowInt(base, exp int64) int64 {
	if exp < 0 {
		return 0
	}
	var result int64 = 1
	for exp > 0 {
		if exp&1 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}
	return result
}

// IsPrime reports whether n is a prime number.
func IsPrime(n int64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := int64(5); i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// PrimesBelow returns all primes strictly less than n.
func PrimesBelow(n int) []int {
	if n <= 2 {
		return nil
	}
	sieve := make([]bool, n)
	for i := 2; i < n; i++ {
		sieve[i] = true
	}
	for i := 2; i*i < n; i++ {
		if !sieve[i] {
			continue
		}
		for j := i * i; j < n; j += i {
			sieve[j] = false
		}
	}
	out := make([]int, 0, n/10)
	for i := 2; i < n; i++ {
		if sieve[i] {
			out = append(out, i)
		}
	}
	return out
}

// Fibonacci returns the n-th Fibonacci number (0-indexed: F(0)=0, F(1)=1).
// n < 0 returns 0.
func Fibonacci(n int) int64 {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	var a, b int64 = 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
