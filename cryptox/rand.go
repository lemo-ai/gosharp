package cryptox

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// RandomBytes returns n cryptographically secure random bytes.
func RandomBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidKeySize
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}

// RandomHex returns a hex string of n random bytes (length 2*n).
func RandomHex(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
