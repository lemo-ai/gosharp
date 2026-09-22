package cryptox

import (
	"bytes"
	"testing"
)

func TestAESGCMRoundTrip(t *testing.T) {
	key, err := GenerateAESKey(32)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("hello gosharp cryptox")
	sealed, err := AESGCMEncrypt(key, plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := AESGCMDecrypt(key, sealed)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q want %q", got, plain)
	}

	b64, err := AESGCMEncryptBase64(key, string(plain))
	if err != nil {
		t.Fatal(err)
	}
	s, err := AESGCMDecryptBase64(key, b64)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(plain) {
		t.Fatalf("base64 roundtrip: got %q", s)
	}
}

func TestAESGCMBadKey(t *testing.T) {
	if _, err := AESGCMEncrypt([]byte("short"), []byte("x")); err != ErrInvalidKeySize {
		t.Fatalf("want ErrInvalidKeySize, got %v", err)
	}
}

func TestChaCha20RoundTrip(t *testing.T) {
	key, err := GenerateChaCha20Key()
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("chacha20-poly1305")
	sealed, err := ChaCha20Encrypt(key, plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ChaCha20Decrypt(key, sealed)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q want %q", got, plain)
	}
}

func TestRSAEncryptSign(t *testing.T) {
	priv, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("sign me")
	ct, err := RSAEncrypt(&priv.PublicKey, msg)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := RSADecrypt(priv, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, msg) {
		t.Fatalf("rsa decrypt mismatch")
	}

	sig, err := RSASign(priv, msg)
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerify(&priv.PublicKey, msg, sig); err != nil {
		t.Fatal(err)
	}
	if err := RSAVerify(&priv.PublicKey, []byte("tampered"), sig); err != ErrVerifyFailed {
		t.Fatalf("want ErrVerifyFailed, got %v", err)
	}

	privPEM, err := MarshalPrivateKeyPEM(priv)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM, err := MarshalPublicKeyPEM(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	priv2, err := ParsePrivateKeyPEM(privPEM)
	if err != nil {
		t.Fatal(err)
	}
	pub2, err := ParsePublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatal(err)
	}
	ct2, err := RSAEncrypt(pub2, msg)
	if err != nil {
		t.Fatal(err)
	}
	pt2, err := RSADecrypt(priv2, ct2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt2, msg) {
		t.Fatal("pem roundtrip failed")
	}
}

func TestHMAC(t *testing.T) {
	key := []byte("secret")
	data := []byte("payload")
	mac := HMACSHA256(key, data)
	if !HMACSHA256Equal(key, data, mac) {
		t.Fatal("equal failed")
	}
	if HMACSHA256Equal(key, []byte("other"), mac) {
		t.Fatal("should not equal")
	}
	hex := HMACSHA256Hex(key, data)
	if len(hex) != 64 {
		t.Fatalf("hex len=%d", len(hex))
	}
}

func TestPassword(t *testing.T) {
	hash, err := HashPasswordCost("s3cret!", 4) // low cost for tests
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "s3cret!") {
		t.Fatal("check failed")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("should mismatch")
	}
	if err := CheckPasswordErr(hash, "wrong"); err != ErrPasswordMismatch {
		t.Fatalf("want ErrPasswordMismatch, got %v", err)
	}

	salt, err := RandomBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	k1 := PBKDF2Key([]byte("pass"), salt, 1000, 32)
	k2 := PBKDF2Key([]byte("pass"), salt, 1000, 32)
	if !bytes.Equal(k1, k2) {
		t.Fatal("pbkdf2 not deterministic")
	}
	if len(k1) != 32 {
		t.Fatalf("key len=%d", len(k1))
	}
}

func TestRandom(t *testing.T) {
	a, err := RandomBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	b, err := RandomBytes(16)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("unlikely equal random")
	}
	h, err := RandomHex(8)
	if err != nil {
		t.Fatal(err)
	}
	if len(h) != 16 {
		t.Fatalf("hex len=%d", len(h))
	}
}
