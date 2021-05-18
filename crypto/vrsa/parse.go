package vrsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// ParsePKIXPublicKey 解析rsa pkix publicKey
func ParsePKIXPublicKey(pub io.Reader) (pk *rsa.PublicKey, err error) {
	publicKey, err := io.ReadAll(pub)
	if err != nil {
		return
	}
	block, _ := pem.Decode(publicKey)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pk = pubInterface.(*rsa.PublicKey)
	return
}

// ParsePKCS1PrivateKey 解析rsa pkcs1 privateKey
func ParsePKCS1PrivateKey(priv io.Reader) (priKey *rsa.PrivateKey, err error) {
	privateKey, err := io.ReadAll(priv)
	if err != nil {
		return
	}
	block, _ := pem.Decode(privateKey)
	priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

// ParsePKCS8PrivateKey 解析rsa pkcs8 privateKey
func ParsePKCS8PrivateKey(priv io.Reader) (privKey *rsa.PrivateKey, err error) {
	privateKey, err := io.ReadAll(priv)
	if err != nil {
		return
	}
	block, _ := pem.Decode(privateKey)
	privKeyInt, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	privKey = privKeyInt.(*rsa.PrivateKey)
	return
}
