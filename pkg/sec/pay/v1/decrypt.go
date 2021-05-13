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

/*
解密流程

step1 将密文块分包
	各部分数据,及主体数据
	1.支付密文 2.aes密钥密文 3.客户端rsa公钥密文 4.签名密文
	body密文 = 支付密文.aes密钥.客户端rsa公钥
	payCipherText, aesKeyCipherText, rsaPubCipherText, sig, bodyCipherText := payCipherBlock(payCipher)
step2 验证主体数据签名
	_, err = decryptVerifySign(rsaPubCipherText, bodyCipherText, sig)
	if err != nil {
		return
	}
step3 解密aes密码
step4 解密支付密文
*/

// DecryptWithFn 通过钩子函数读取密钥并解密
func DecryptWithFn(payCipher string, pf func(pubStr string) (privStr string, err error)) (payData string, err error) {
	payCipherText, aesKeyCipherText, rsaPubCipherText, sig, bodyCipherText := payCipherBlock(payCipher)
	pubKey, err := decryptVerifySign(rsaPubCipherText, bodyCipherText, sig)
	if err != nil {
		return
	}
	privStr, err := pf(pubKey)
	if err != nil {
		return
	}
	payData, err = decryptCipher(privStr, aesKeyCipherText, payCipherText)
	return
}

// Decrypt 服务端支付解密
func Decrypt(payCipher string, privStr string) (payData string, err error) {
	// 解密流程

	// step1 将密文块分包
	// 各部分数据,及主体数据
	// 1.支付密文 2.aes密钥密文 3.客户端rsa公钥密文 4.签名密文
	// body密文 = 支付密文.aes密钥.客户端rsa公钥
	payCipherText, aesKeyCipherText, rsaPubCipherText, sig, bodyCipherText := payCipherBlock(payCipher)
	// step2 验证主体数据签名
	_, err = decryptVerifySign(rsaPubCipherText, bodyCipherText, sig)
	if err != nil {
		return
	}
	// step3 解密aes密码
	// step4 解密支付密文
	payData, err = decryptCipher(privStr, aesKeyCipherText, payCipherText)
	return
}

// decryptVerifySign 解密验签
func decryptVerifySign(rsaPubCipherText, bodyCipherText, sig string) (pubStr string, err error) {
	// 验证主体数据签名
	digest, err := vsha.Sha256Hash(bodyCipherText)
	if err != nil {
		return
	}
	rsaPub, err := base64.StdEncoding.DecodeString(rsaPubCipherText)
	if err != nil {
		return
	}
	var pubBuf = bytes.NewBuffer(nil)
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
	pubStr = string(rsaPub)
	return
}

// decryptCipher 解密密文
func decryptCipher(privStr, aesKeyCipherText, payCipherText string) (payData string, err error) {
	// step3 解密aes密码
	var serverPrivBuf = bytes.NewBuffer(nil)
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

// payCipherBlock 拆分密文块
func payCipherBlock(payCipher string) (payCipherText, aesKeyCipherText, rsaPubCipherText, sig, bodyCipherText string) {
	// step1 将密文块分包
	// 各部分数据,及主体数据
	// 1.支付密文 2.aes密钥密文 3.客户端rsa公钥密文 4.签名密文
	// body密文 = 支付密文.aes密钥.客户端rsa公钥
	cps := strings.Split(payCipher, ".")
	payCipherText = cps[0]
	aesKeyCipherText = cps[1]
	rsaPubCipherText = cps[2]
	sig = cps[3]
	bodyCipherText = payCipherText + "." + aesKeyCipherText + "." + rsaPubCipherText
	return
}
