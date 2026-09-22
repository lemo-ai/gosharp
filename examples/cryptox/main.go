package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/cryptox"
)

func main() {
	key, err := cryptox.GenerateAESKey(32)
	if err != nil {
		panic(err)
	}
	cipherB64, err := cryptox.AESGCMEncryptBase64(key, "hello cryptox")
	if err != nil {
		panic(err)
	}
	plain, err := cryptox.AESGCMDecryptBase64(key, cipherB64)
	if err != nil {
		panic(err)
	}
	fmt.Println("AES-GCM:", plain)

	priv, err := cryptox.GenerateRSAKey(2048)
	if err != nil {
		panic(err)
	}
	msg := []byte("signed message")
	sig, err := cryptox.RSASign(priv, msg)
	if err != nil {
		panic(err)
	}
	if err := cryptox.RSAVerify(&priv.PublicKey, msg, sig); err != nil {
		panic(err)
	}
	fmt.Println("RSA-PSS: ok")

	mac := cryptox.HMACSHA256Hex([]byte("key"), []byte("data"))
	fmt.Println("HMAC:", mac[:16]+"...")

	hash, err := cryptox.HashPassword("password")
	if err != nil {
		panic(err)
	}
	fmt.Println("bcrypt match:", cryptox.CheckPassword(hash, "password"))
}
