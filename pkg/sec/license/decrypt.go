package license

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"

	"github.com/vicnoah/golib/pkg/sec/vaes"
	"github.com/vicnoah/golib/pkg/sec/vrsa"
)

// Decrypt Decrypting gob encrypted data based on aes-cbc-128
func Decrypt(data, salt, rsaPublicKey []byte) (ret []byte, err error) {
	var (
		payload Payload
		buf     bytes.Buffer
	)
	_, err = buf.Write(data)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&payload)
	if err != nil {
		return
	}
	passHash, err := vrsa.PubKeyDecrypt(rsaPublicKey, payload.Rsa)
	if err != nil {
		return
	}
	h1 := sha1.New()
	h1.Write(salt)
	saltHash := h1.Sum(nil)
	passwd := append(passHash, saltHash...)
	h2 := sha1.New()
	h2.Write(passwd)
	hashPasswd := h2.Sum(nil)
	ret, err = vaes.CBCDecrypt(payload.Aes, hashPasswd, []byte(_IV))
	return
}
