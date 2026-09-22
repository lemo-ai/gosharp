// Package hashx provides convenience hashing helpers.
package hashx

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

// MD5 returns the hex-encoded MD5 of data.
func MD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// MD5String returns the hex-encoded MD5 of s.
func MD5String(s string) string {
	return MD5([]byte(s))
}

// SHA1 returns the hex-encoded SHA-1 of data.
func SHA1(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

// SHA256 returns the hex-encoded SHA-256 of data.
func SHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// SHA256String returns the hex-encoded SHA-256 of s.
func SHA256String(s string) string {
	return SHA256([]byte(s))
}

// FNV64a returns a non-cryptographic 64-bit FNV-1a hash.
func FNV64a(data []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(data)
	return h.Sum64()
}
