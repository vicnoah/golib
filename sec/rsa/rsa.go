// Package rsa RSA encryption and decryption algorithm implementation
// rsa has no use of private key encryption and public key decryption.
// Private key encryption and public key decryption are implemented based on signature verification.
// The type is determined as a byte slice. When using it, pay attention to the matching of the signature key.
// Related article links https://www.cnblogs.com/imlgc/p/7076313.html
package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	crsa "git.sabertrain.com/vector-tech/golib/sec/crypt/rsa"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

var (
	ErrPublicKey = errors.New("public key error")
	ErrPrivateKey = errors.New("private key error")
	ErrParsePublicKey = errors.New("parse public key error")
)

// PubKeyEncrypt Public key encryption
// publicKey is the public key of pkcs1
func PubKeyEncrypt(publicKey, originData []byte) (cipher []byte, err error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		err = ErrPublicKey
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		err = ErrParsePublicKey
		return
	}
	cipher, err = rsa.EncryptPKCS1v15(rand.Reader, pub, originData)
	return
}

// PriKeyDecrypt Private key decryption
func PriKeyDecrypt(privateKey, cipher []byte) (data []byte, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err = ErrPrivateKey
		return
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	data, err = rsa.DecryptPKCS1v15(rand.Reader, priv, cipher)
	return
}

// PriKeyEncrypt Private key encryption
func PriKeyEncrypt(privateKey, data []byte) (cipher []byte, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err = ErrPrivateKey
		return
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	cipher, err =  crsa.PrivateEncrypt(priv, data)
	return
}

// PubKeyDecrypt Public key decryption
func PubKeyDecrypt(publicKey, cipher []byte) (data []byte, err error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		err = ErrPublicKey
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		err = ErrParsePublicKey
		return
	}
	data, err = crsa.PublicDecrypt(pub, cipher)
	return
}

// SignWithSha1Hex Signature based on pkcs1, output hex format
func SignWithSha1Hex(privateKey, data []byte) (sign []byte, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err = ErrPrivateKey
		return
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	h := sha1.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	sign, err = rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, hash[:])
	return
}

// VerifyWithSha1Hex Pkcs1 signature to verify hex data
func VerifyWithSha1Hex(publicKey, data, signData []byte) (err error) {
	sign, err := hex.DecodeString(string(signData))
	if err != nil {
		return
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		err = ErrPublicKey
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		err = ErrParsePublicKey
		return
	}
	hash := sha1.New()
	hash.Write(data)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hash.Sum(nil), sign)
	return
}