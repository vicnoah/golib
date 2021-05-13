package v1

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/vicnoah/golib/crypto/vaes"
	"github.com/vicnoah/golib/crypto/vrsa"
	"github.com/vicnoah/golib/crypto/vsha"
	"github.com/vicnoah/golib/pkg/rand"
)

const (
	payAesKeyLength = 16 // aes密钥长度,必须是16的整数倍
)

var (
	// ErrPayEnvironmentNotSecure 支付环境不安全错误
	ErrPayEnvironmentNotSecure = errors.New("pay environment not secure")
)

// Encrypt 支付客户端完成数据加密
func Encrypt(keyPair string, payload string) (payCipher string, err error) {
	// 支付环境安全验证
	ok, serverPubStr, err := VerifyPss(keyPair)
	if err != nil {
		return
	}
	if !ok {
		err = ErrPayEnvironmentNotSecure
		return
	}
	// 支付数据加解密
	// 处理key pair

	// step1 生成客户端随机rsa pkcs#1加密对
	// 此密钥对用于对数据签名
	var (
		privBuf = bytes.NewBuffer(nil)
		pubBuf  = bytes.NewBuffer(nil)
	)
	// 生成密钥
	err = vrsa.GenerateRSAKey(privBuf, pubBuf, 1024)
	if err != nil {
		return
	}
	// privStr := privBuf.String()
	pubStr := pubBuf.String()

	// step2 生成随机aes加密密钥
	aesKey := rand.GetString(payAesKeyLength)

	// step3 加密aes密钥对
	var serverPubBuf = bytes.NewBuffer(nil)
	_, err = serverPubBuf.Write([]byte(serverPubStr))
	if err != nil {
		return
	}
	aesKeyCipherText, err := vrsa.EncryptOAEP(serverPubBuf, aesKey)
	if err != nil {
		return
	}

	// step4 支付数据加密
	payCipherBs, err := vaes.CBCEncrypt([]byte(payload), []byte(aesKey))
	if err != nil {
		return
	}
	payCipherText := base64.StdEncoding.EncodeToString(payCipherBs)

	// step5 主体密文 = 将加密aes密钥对,支付密文,客户端签名公钥进行连接
	// 主体密文 = 支付密文 + 加密(aes密钥) + 客户端随机rsa公钥 + 服务端rsa公钥
	// 支付密文 = base64(aes256(payload))
	// aes密钥 = base64(RSA OAEP(aes密钥))
	// 客户端公钥 = base64(rsa随机公钥)
	// 服务端公钥 = base64(服务端rsa公钥)
	b64PubStr := base64.StdEncoding.EncodeToString([]byte(pubStr))
	b64ServerPubStr := base64.StdEncoding.EncodeToString([]byte(serverPubStr))
	payCipherArr := make([]string, 0)
	payCipherArr = append(payCipherArr, payCipherText, aesKeyCipherText, b64PubStr, b64ServerPubStr)
	cipherBody := strings.Join(payCipherArr, ".")

	// step6 使用客户端随机rsa私钥进行sha256数据签名并连接主体密文
	// 签名也是签主体密文
	// 最终密文 = 主体密文 + 签名(主体密文)
	// 密钥格式: 采用点分形式 支付密文.aes密钥.客户端rsa公钥.服务端rsa公钥.签名
	// 密文示例: OipsPI=.oWbWKRUU=.6+nEm9wmcT/bW.6+nEm9wmcT/bW.Em9wmcT/bWrchg
	digest, err := vsha.Sha256Hash(cipherBody)
	if err != nil {
		return
	}
	privateKey, err := vrsa.ParsePKCS1PrivateKey(privBuf)
	if err != nil {
		return
	}
	sigBody, err := vrsa.SignPSS(privateKey, digest)
	if err != nil {
		return
	}
	payCipherArr = append(payCipherArr, sigBody)
	payCipher = strings.Join(payCipherArr, ".")
	return
}

// VerifyPss 支付客户端验证密钥签名
func VerifyPss(keyPair string) (ok bool, pubStr string, err error) {
	// base64解密密钥对
	kp, err := base64.StdEncoding.DecodeString(keyPair)
	if err != nil {
		fmt.Println(err)
	}
	// 截取密钥对中公钥及签名部分
	kps := strings.Split(string(kp), ".")
	pubStr = kps[0]
	sig := kps[1]
	// 使用sha256计算公钥签名
	digest, err := vsha.Sha256Hash(kps[0])
	if err != nil {
		return
	}
	// 读取并使用公钥
	buf := bytes.NewBuffer(nil)
	buf.WriteString(pubStr)
	pubKey, err := vrsa.ParsePKIXPublicKey(buf)
	if err != nil {
		return
	}
	// 使用公钥,公钥摘要,sig签名进行验签
	if vrsa.VerifyPSS(pubKey, digest, sig) {
		ok = true
		return
	}
	return
}
