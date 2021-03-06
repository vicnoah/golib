package v1

import (
	"bytes"
	"encoding/base64"

	"github.com/vicnoah/golib/crypto/vrsa"
	"github.com/vicnoah/golib/crypto/vsha"
)

// GrabKey 支付方抓取密钥
// keyPair: 密钥对 = base64(密钥.签名<rsapss(sha256(公钥))>)
func GrabKey(bits int) (privStr, pubStr, keyPair string, err error) {
	var (
		privBuf = bytes.NewBuffer(nil)
		pubBuf  = bytes.NewBuffer(nil)
	)
	// 生成密钥
	err = vrsa.GenerateRSAKey(privBuf, pubBuf, bits)
	if err != nil {
		return
	}
	privStr = privBuf.String()
	pubStr = pubBuf.String()
	// 生成公有密钥签名,使用sha256生成摘要
	privateKey, err := vrsa.ParsePKCS8PrivateKey(privBuf)
	if err != nil {
		return
	}
	digest, err := vsha.Sha256Hash(pubStr)
	if err != nil {
		return
	}
	// 对公钥进行签名
	sig, err := vrsa.SignPSS(privateKey, digest)
	if err != nil {
		return
	}
	// 连接公钥密钥对
	keyPair = base64.StdEncoding.EncodeToString([]byte(pubStr + "." + sig))
	return
}
