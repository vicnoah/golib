package v1

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/vicnoah/golib/crypto/vaes"
	"github.com/vicnoah/golib/crypto/vrsa"
	"github.com/vicnoah/golib/crypto/vsha"
)

var (
	// ErrPayDataTamper 支付数据被篡改
	ErrPayDataTamper = errors.New("pay data is tampered")
)

// Decrypt 服务端支付解密
func Decrypt(payCipher string, privStr string) (payData string, err error) {
	// 解密流程

	// step1 将密文块分包
	// 各部分数据,及主体数据
	// 1.支付密文 2.aes密钥密文 3.客户端rsa公钥密文 4.签名密文
	// body密文 = 支付密文.aes密钥.客户端rsa公钥
	cps := strings.Split(payCipher, ".")
	payCipherText := cps[0]
	aesKeyCipherText := cps[1]
	rsaPubCipherText := cps[2]
	sig := cps[3]
	bodyCipherText := payCipherText + "." + aesKeyCipherText + "." + rsaPubCipherText
	// step2 验证主体数据签名
	digest, err := vsha.Sha256Hash(bodyCipherText)
	if err != nil {
		return
	}
	rsaPub, err := base64.StdEncoding.DecodeString(rsaPubCipherText)
	if err != nil {
		return
	}
	var pubBuf = bytes.NewBuffer([]byte{})
	_, err = pubBuf.Write(rsaPub)
	if err != nil {
		return
	}
	pubKey, err := vrsa.ParsePKIXPublicKey(pubBuf)
	if err != nil {
		return
	}
	// 数据签名验证
	if !vrsa.VerifyPSS(pubKey, digest, sig) {
		err = ErrPayDataTamper
		return
	}
	// step3 解密aes密码
	var serverPrivBuf = bytes.NewBuffer([]byte{})
	_, err = serverPrivBuf.Write([]byte(privStr))
	if err != nil {
		return
	}
	aesKey, err := vrsa.DecryptOAEP(serverPrivBuf, aesKeyCipherText)
	if err != nil {
		return
	}
	// step4 解密支付密文
	payCipherBs, err := base64.StdEncoding.DecodeString(payCipherText)
	if err != nil {
		return
	}
	payDataBs, err := vaes.CBCDecrypt(payCipherBs, []byte(aesKey))
	if err != nil {
		return
	}
	payData = string(payDataBs)
	return
}
