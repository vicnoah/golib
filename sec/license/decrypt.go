package license

import (
	"bytes"
	"crypto/sha1"
	"e.coding.net/vector-tech/golib/sec/aes"
	"e.coding.net/vector-tech/golib/sec/rsa"
	"encoding/gob"
)
// Decrypt Decrypting gob encrypted data based on aes-cbc-128
func Decrypt(data, salt, rsaPublicKey []byte) (ret []byte, err error) {
	var (
		payload Payload
		buf bytes.Buffer
		)
	_, err = buf.Write(data)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(payload)
	if err != nil {
		return
	}
	passHash, err := rsa.PubKeyDecrypt(rsaPublicKey, payload.Rsa)
	if err != nil {
		return
	}
	h1 := sha1.New()
	h1.Write(salt)
	saltHash := h1.Sum(nil)
	passwd := append(passHash, saltHash...)
	ret, err = aes.CBCDecrypt(payload.Aes, passwd, []byte(_IV))
	return
}
