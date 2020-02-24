package license

import (
	"bytes"
	"crypto/sha1"
	"e.coding.net/vector-tech/golib/sec/aes"
	"e.coding.net/vector-tech/golib/sec/rsa"
	"encoding/gob"
)

const (
	_IV = "vector"
)

// Encrypt Gob-encoded key generated based on aes-cbc-128
func Encrypt(data, pass, salt, rsaPrivateKey []byte) (ret []byte, err error) {
	var payload Payload
	// 取hash值
	h := sha1.New()
	h.Write(pass)
	passHash := h.Sum(nil)
	h1 := sha1.New()
	h1.Write(salt)
	saltHash := h1.Sum(nil)
	passwd := append(passHash, saltHash...)
	h2 := sha1.New()
	h2.Write(passwd)
	hashPasswd := h2.Sum(nil)
	payload.Aes, err = aes.CBCEncrypt(data, hashPasswd, []byte(_IV))
	if err != nil {
		return
	}
	payload.Rsa, err = rsa.PriKeyEncrypt(rsaPrivateKey, passHash)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(payload)
	ret = buf.Bytes()
	return
}