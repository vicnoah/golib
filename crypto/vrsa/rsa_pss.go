package vrsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

// SignPSS rsa pss签名算法
// signature为base64编码
func SignPSS(priv *rsa.PrivateKey, digest string) (signature string, err error) {
	signatureBs, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, []byte(digest), nil)
	signature = base64.StdEncoding.EncodeToString(signatureBs)
	return
}

// VerifyPSS rsa pss验签算法
func VerifyPSS(pub *rsa.PublicKey, digest, signature string) (ok bool) {
	bs, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return
	}
	err = rsa.VerifyPSS(pub, crypto.SHA256, []byte(digest), bs, nil)
	if err == nil {
		ok = true
		return
	}
	return
}
