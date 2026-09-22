// Package randx provides random helpers for tests and lightweight IDs.
package randx

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	mrand "math/rand/v2"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Intn returns a crypto-strength random int in [0, n). n must be > 0.
func Intn(n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}

// String returns a random alphanumeric string of length n.
func String(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b := make([]byte, n)
	max := big.NewInt(int64(len(letters)))
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = letters[v.Int64()]
	}
	return string(b), nil
}

// Hex returns n random bytes encoded as hex (2n chars).
func Hex(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// FastIntn is a non-crypto PRNG helper for non-security use cases.
func FastIntn(n int) int {
	if n <= 0 {
		return 0
	}
	return mrand.IntN(n)
}

// Uint64 returns a crypto-strength random uint64.
func Uint64() (uint64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}
