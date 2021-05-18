package v1

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/vicnoah/golib/crypto/vaes"
	"github.com/vicnoah/golib/crypto/vrsa"
	"github.com/vicnoah/golib/crypto/vsha"
	"github.com/vicnoah/golib/pkg/rand"
)

const (
	aesKeyLength = 16 // aes密钥长度,必须是16的整数倍
)

func TestPay(t *testing.T) {
	// 为支付app抓取密钥
	privStr, pubStr, keyPair, err := GrabKey(1024)
	_ = pubStr
	if err != nil {
		fmt.Printf("server端密钥生成错误: %v\n", err)
		return
	}
	// app端进行支付加密
	payCipher, err := Encrypt(keyPair, `
{"password": "sfdjsdf"}	
	`)
	if err != nil {
		fmt.Printf("支付失败: %v\n", err)
		return
	}
	// server得到支付密文
	// server端解密支付密文
	payData, err := Decrypt(payCipher, privStr)
	if err != nil {
		fmt.Printf("支付失败: %v\n", err)
		return
	}
	fmt.Println(payData)
	fmt.Println("支付成功")
}

func TestClientPay(t *testing.T) {
	keyPair := `LS0tLS1CRUdJTiBSU0EgUHVibGljIEtleS0tLS0tCk1JR2ZNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0R05BRENCaVFLQmdRREhzelpjN3pGUWR1NmZpSFVkdmszaW1CT08KbGFnU0xrVFNSdWQ0dFpUcWdpbTV4eDlsUlpmajloOHkzZjJFeWNBcENKOGo4WWNZVjRsZGxVY09yZlYyQWk4LwpaeDhNRW9tTFFQWkM2TWc3dHNrSlpVL2FIZVpnUXRueklIS1FEdC9EdEh1Zk5CNlhzcGVmZGlVVUhPNFovM3dDCmVZZmVyMEZTdHhCN1Y0TXdWUUlEQVFBQgotLS0tLUVORCBSU0EgUHVibGljIEtleS0tLS0tCi5ka3Y5ejdPRU1oWDZUWmNIcStTVXYvd0ZvUVVtaEt0TzlBU202YnZwQkpUeFFBdlNpbEFMMk5RZzdNaU01VURPdjRkSlFqMG4waUdRMTBGWXlpblBsUkNnS1dOZlFhQkFNRktUUlE1a1Y3aWpsZ20wVnRxUC9PVDNBQjhKbC93blYrMWdDS0NBUWpYME1aRlJyRzkvaGtkUUxNdXpPSjNyUVBqN0VmbC9keVU9`
	// 支付环境安全验证
	ok, serverPubStr, err := clientVerifyPss(keyPair)
	if err != nil {
		return
	}
	if !ok {
		fmt.Println("签名验证失败,支付环境不安全")
		err = errors.New("pay environment not secure")
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
	aesKey := rand.GetString(aesKeyLength)

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
	payload := `
{
	"password": "404ae203e933c1f70037f2450e77a2a5",
	"captcha": "739453"
}
	`
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
	serverB64PubStr := base64.StdEncoding.EncodeToString([]byte(serverPubStr))
	payCipherArrs := make([]string, 0)
	payCipherArrs = append(payCipherArrs, payCipherText, aesKeyCipherText, b64PubStr, serverB64PubStr)
	cipherBody := strings.Join(payCipherArrs, ".")

	// step6 使用客户端随机rsa私钥进行sha256数据签名并连接主体密文
	// 签名也是签主体密文
	// 最终密文 = 主体密文 + 签名(主体密文)
	// 密钥格式: 采用点分形式 支付密文.aes密钥.客户端rsa公钥.服务端rsa公钥.签名
	// 密文示例: OipsPI=.oWbWKRUU=.6+nEm9wmcT/bW.6+nEm9wmcT/bW.Em9wmcT/bWrchg
	digest, err := vsha.Sha256Hash(cipherBody)
	if err != nil {
		return
	}
	privateKey, err := vrsa.ParsePKCS8PrivateKey(privBuf)
	if err != nil {
		return
	}
	sigBody, err := vrsa.SignPSS(privateKey, digest)
	if err != nil {
		return
	}
	// 密文加签名
	payCipherArrs = append(payCipherArrs, sigBody)
	payCipher := strings.Join(payCipherArrs, ".")
	fmt.Println(payCipher)
	return
}

func TestPayCheck(t *testing.T) {
	// 为支付app抓取密钥
	privStr, pubStr, keyPair, err := GrabKey(1024)
	_ = pubStr
	if err != nil {
		fmt.Printf("server端密钥生成错误: %v\n", err)
		return
	}
	// app端进行支付加密
	payCipher, err := clientPay(keyPair)
	if err != nil {
		fmt.Printf("app端支付加密错误: %v", err)
		return
	}
	// server得到支付密文
	// server端解密支付密文
	payData, err := serverPay(payCipher, privStr)
	if err != nil {
		fmt.Printf("支付失败: %v\n", err)
		return
	}
	fmt.Println(payData)
	fmt.Println("支付成功")
}

// serverPay 服务端支付验证
func serverPay(payCipher string, privStr string) (payData string, err error) {
	// 解密流程

	// step1 将密文块分包
	// 各部分数据,及主体数据
	// 1.支付密文 2.aes密钥密文 3.客户端rsa公钥密文 4.服务端rsa公钥密文 5.签名密文
	// body密文 = 支付密文.aes密钥.客户端rsa公钥.服务端rsa公钥
	cps := strings.Split(payCipher, ".")
	payCipherText := cps[0]
	aesKeyCipherText := cps[1]
	rsaPubCipherText := cps[2]
	serverRsaPubCipherText := cps[3]
	sig := cps[4]
	// 生成主体密文
	bodyCipherArr := make([]string, 0)
	bodyCipherArr = append(bodyCipherArr, payCipherText, aesKeyCipherText, rsaPubCipherText, serverRsaPubCipherText)
	bodyCipherText := strings.Join(bodyCipherArr, ".")
	// step2 验证主体数据签名
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
		fmt.Println("签名验证失败,数据被篡改")
		return
	}
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

// clientPay 客户端支付
func clientPay(keyPair string) (payCipher string, err error) {
	// 支付环境安全验证
	ok, serverPubStr, err := clientVerifyPss(keyPair)
	if err != nil {
		return
	}
	if !ok {
		fmt.Println("签名验证失败,支付环境不安全")
		err = errors.New("pay environment not secure")
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
	aesKey := rand.GetString(aesKeyLength)

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
	payload := `
{
	"pay_pass": "sdjfsdjfgjsdfjsdj",
	"cost": 107.00
}
	`
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
	serverB64PubStr := base64.StdEncoding.EncodeToString([]byte(serverPubStr))
	payCipherArrs := make([]string, 0)
	payCipherArrs = append(payCipherArrs, payCipherText, aesKeyCipherText, b64PubStr, serverB64PubStr)
	cipherBody := strings.Join(payCipherArrs, ".")

	// step6 使用客户端随机rsa私钥进行sha256数据签名并连接主体密文
	// 签名也是签主体密文
	// 最终密文 = 主体密文 + 签名(主体密文)
	// 密钥格式: 采用点分形式 支付密文.aes密钥.客户端rsa公钥.服务端rsa公钥.签名
	// 密文示例: OipsPI=.oWbWKRUU=.6+nEm9wmcT/bW.Em9wmcT/bWrchg
	digest, err := vsha.Sha256Hash(cipherBody)
	if err != nil {
		return
	}
	privateKey, err := vrsa.ParsePKCS8PrivateKey(privBuf)
	if err != nil {
		return
	}
	sigBody, err := vrsa.SignPSS(privateKey, digest)
	if err != nil {
		return
	}
	// 密文加签名
	payCipherArrs = append(payCipherArrs, sigBody)
	payCipher = strings.Join(payCipherArrs, ".")
	return
}

// clientVerifyPss 客户端验证密钥签名
func clientVerifyPss(keyPair string) (ok bool, pubStr string, err error) {
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
