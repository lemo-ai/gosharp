package json

import "strconv"

// parseIntFast parses a decimal integer from b without allocations.
// Returns ok=false if the token is not a plain integer (e.g. has '.'/'e').
func parseIntFast(b []byte) (int64, bool) {
	if len(b) == 0 {
		return 0, false
	}
	i := 0
	neg := false
	switch b[0] {
	case '-':
		neg = true
		i++
		if i >= len(b) {
			return 0, false
		}
	case '+':
		i++
		if i >= len(b) {
			return 0, false
		}
	}
	var n uint64
	for ; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + uint64(c-'0')
	}
	if neg {
		return -int64(n), true
	}
	return int64(n), true
}

func parseFloatBytes(b []byte) (float64, error) {
	return strconv.ParseFloat(b2s(b), 64)
}

func appendInt(dst []byte, n int64) []byte {
	return strconv.AppendInt(dst, n, 10)
}
