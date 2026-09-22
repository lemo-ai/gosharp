package cryptox

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HMACSHA256 returns the raw HMAC-SHA256 of data.
func HMACSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}

// HMACSHA256Hex returns the hex-encoded HMAC-SHA256 of data.
func HMACSHA256Hex(key, data []byte) string {
	return hex.EncodeToString(HMACSHA256(key, data))
}

// HMACSHA256Equal reports whether mac equals HMAC-SHA256(key, data)
// using constant-time comparison.
func HMACSHA256Equal(key, data, mac []byte) bool {
	return hmac.Equal(HMACSHA256(key, data), mac)
}
