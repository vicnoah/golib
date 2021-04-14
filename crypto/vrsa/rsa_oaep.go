package vrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"io"
)

// EncryptOAEP rsa OAEP加密算法
func EncryptOAEP(pub io.Reader, text string) (ciphertext string, err error) {
	rsaPublicKey, err := ParsePKIXPublicKey(pub)
	if err != nil {
		return
	}
	secretMessage := []byte(text)
	rng := rand.Reader
	cipherdata, err := rsa.EncryptOAEP(sha1.New(), rng, rsaPublicKey, secretMessage, nil)
	if err != nil {
		return
	}
	ciphertext = base64.StdEncoding.EncodeToString(cipherdata)
	return
}

// DecryptOAEP rsa OAEP解密算法
func DecryptOAEP(priv io.Reader, ciphertext string) (text string, err error) {
	rsaPrivateKey, err := ParsePKCS1PrivateKey(priv)
	if err != nil {
		return
	}
	cipherdata, _ := base64.StdEncoding.DecodeString(ciphertext)
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha1.New(), rng, rsaPrivateKey, cipherdata, nil)
	if err != nil {
		return
	}
	text = string(plaintext)
	return
}
