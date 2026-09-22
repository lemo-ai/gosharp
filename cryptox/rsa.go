package cryptox

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const defaultRSABits = 2048

// GenerateRSAKey generates an RSA private key. bits defaults to 2048 when <= 0.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits <= 0 {
		bits = defaultRSABits
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// RSAEncrypt encrypts plaintext with RSA-OAEP (SHA-256).
func RSAEncrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	if pub == nil {
		return nil, ErrNilKey
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, nil)
}

// RSADecrypt decrypts ciphertext produced by RSAEncrypt.
func RSADecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	if priv == nil {
		return nil, ErrNilKey
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
}

// RSASign signs message with RSA-PSS (SHA-256).
func RSASign(priv *rsa.PrivateKey, message []byte) ([]byte, error) {
	if priv == nil {
		return nil, ErrNilKey
	}
	sum := sha256.Sum256(message)
	return rsa.SignPSS(rand.Reader, priv, crypto.SHA256, sum[:], nil)
}

// RSAVerify verifies an RSA-PSS signature. Returns nil on success.
func RSAVerify(pub *rsa.PublicKey, message, signature []byte) error {
	if pub == nil {
		return ErrNilKey
	}
	sum := sha256.Sum256(message)
	if err := rsa.VerifyPSS(pub, crypto.SHA256, sum[:], signature, nil); err != nil {
		return ErrVerifyFailed
	}
	return nil
}

// MarshalPrivateKeyPEM encodes a private key as PKCS#1 PEM.
func MarshalPrivateKeyPEM(priv *rsa.PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, ErrNilKey
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	return pem.EncodeToMemory(block), nil
}

// ParsePrivateKeyPEM parses a PKCS#1 or PKCS#8 RSA private key PEM.
func ParsePrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("cryptox: failed to decode PEM block")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("cryptox: not an RSA private key")
	}
	return priv, nil
}

// MarshalPublicKeyPEM encodes a public key as PKIX PEM.
func MarshalPublicKeyPEM(pub *rsa.PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, ErrNilKey
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: der}
	return pem.EncodeToMemory(block), nil
}

// ParsePublicKeyPEM parses a PKIX RSA public key PEM.
func ParsePublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("cryptox: failed to decode PEM block")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("cryptox: not an RSA public key")
	}
	return pub, nil
}
