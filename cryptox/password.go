package cryptox

import (
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// DefaultBcryptCost is a reasonable default for bcrypt.
	DefaultBcryptCost = bcrypt.DefaultCost
	// DefaultPBKDF2Iter is a reasonable iteration count for PBKDF2-SHA256.
	DefaultPBKDF2Iter = 600_000
)

// HashPassword hashes password with bcrypt (DefaultBcryptCost).
func HashPassword(password string) (string, error) {
	return HashPasswordCost(password, DefaultBcryptCost)
}

// HashPasswordCost hashes password with the given bcrypt cost.
func HashPasswordCost(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword reports whether password matches a bcrypt hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// CheckPasswordErr is like CheckPassword but returns ErrPasswordMismatch on failure.
func CheckPasswordErr(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// PBKDF2Key derives a key from password using PBKDF2-SHA256.
// If iter <= 0, DefaultPBKDF2Iter is used. If keyLen <= 0, 32 is used.
func PBKDF2Key(password, salt []byte, iter, keyLen int) []byte {
	if iter <= 0 {
		iter = DefaultPBKDF2Iter
	}
	if keyLen <= 0 {
		keyLen = 32
	}
	return pbkdf2.Key(password, salt, iter, keyLen, sha256.New)
}
