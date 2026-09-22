package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// AESGCMEncrypt encrypts plaintext with AES-GCM.
// key must be 16, 24, or 32 bytes (AES-128/192/256).
// The returned slice is nonce || ciphertext || tag.
func AESGCMEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeySize
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// AESGCMDecrypt decrypts data produced by AESGCMEncrypt.
func AESGCMDecrypt(key, sealed []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeySize
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(sealed) < nonceSize {
		return nil, ErrCiphertextTooShort
	}
	nonce, ciphertext := sealed[:nonceSize], sealed[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AESGCMEncryptBase64 encrypts plaintext and returns a standard Base64 string.
func AESGCMEncryptBase64(key []byte, plaintext string) (string, error) {
	sealed, err := AESGCMEncrypt(key, []byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// AESGCMDecryptBase64 decrypts a Base64 string produced by AESGCMEncryptBase64.
func AESGCMDecryptBase64(key []byte, sealedB64 string) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(sealedB64)
	if err != nil {
		return "", err
	}
	plain, err := AESGCMDecrypt(key, sealed)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// GenerateAESKey returns a cryptographically secure AES key of the given size
// (16, 24, or 32). Prefer 32 for AES-256.
func GenerateAESKey(size int) ([]byte, error) {
	if size != 16 && size != 24 && size != 32 {
		return nil, ErrInvalidKeySize
	}
	return RandomBytes(size)
}
