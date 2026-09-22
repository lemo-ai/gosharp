package cryptox

import "errors"

var (
	// ErrInvalidKeySize indicates the key length is not valid for the algorithm.
	ErrInvalidKeySize = errors.New("cryptox: invalid key size")
	// ErrCiphertextTooShort indicates the sealed blob is shorter than nonce+tag.
	ErrCiphertextTooShort = errors.New("cryptox: ciphertext too short")
	// ErrNilKey indicates a nil RSA key was provided.
	ErrNilKey = errors.New("cryptox: nil key")
	// ErrVerifyFailed indicates signature verification failed.
	ErrVerifyFailed = errors.New("cryptox: signature verification failed")
	// ErrPasswordMismatch indicates bcrypt comparison failed.
	ErrPasswordMismatch = errors.New("cryptox: password mismatch")
)
