package cryptox

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// ChaCha20Encrypt encrypts plaintext with ChaCha20-Poly1305.
// key must be exactly 32 bytes.
// The returned slice is nonce || ciphertext || tag.
func ChaCha20Encrypt(key, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, ErrInvalidKeySize
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return aead.Seal(nonce, nonce, plaintext, nil), nil
}

// ChaCha20Decrypt decrypts data produced by ChaCha20Encrypt.
func ChaCha20Decrypt(key, sealed []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, ErrInvalidKeySize
	}
	nonceSize := aead.NonceSize()
	if len(sealed) < nonceSize {
		return nil, ErrCiphertextTooShort
	}
	nonce, ciphertext := sealed[:nonceSize], sealed[nonceSize:]
	return aead.Open(nil, nonce, ciphertext, nil)
}

// ChaCha20EncryptBase64 encrypts and returns a standard Base64 string.
func ChaCha20EncryptBase64(key []byte, plaintext string) (string, error) {
	sealed, err := ChaCha20Encrypt(key, []byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// ChaCha20DecryptBase64 decrypts a Base64 string from ChaCha20EncryptBase64.
func ChaCha20DecryptBase64(key []byte, sealedB64 string) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(sealedB64)
	if err != nil {
		return "", err
	}
	plain, err := ChaCha20Decrypt(key, sealed)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// GenerateChaCha20Key returns a 32-byte key for ChaCha20-Poly1305.
func GenerateChaCha20Key() ([]byte, error) {
	return RandomBytes(chacha20poly1305.KeySize)
}
