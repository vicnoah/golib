package vrsa

import (
	"fmt"
	"os"
	"testing"

	"git.sabertrain.com/vector-dev/golib/crypto/vsha"
)

func TestRsa(t *testing.T) {
	// 测试签名
	privFile, _ := os.OpenFile("public.pem", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	pubFile, _ := os.OpenFile("private.pem", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	// 生成密钥
	err := GenerateRSAKey(privFile, pubFile, 1024)
	if err != nil {
		return
	}
	privFile.Close()
	pubFile.Close()
	// 打开文件
	privFile, _ = os.Open("public.pem")
	pubFile, _ = os.Open("private.pem")
	priv, err := ParsePKCS1PrivateKey(privFile)
	if err != nil {
		fmt.Println(err)
	}
	pub, err := ParsePKIXPublicKey(pubFile)
	if err != nil {
		fmt.Println(err)
	}
	text := "hello world"
	hash, _ := vsha.Sha256Hash(text)
	signature, _ := SignPSS(priv, hash)
	if VerifyPSS(pub, hash, signature) {
		fmt.Println("签名通过")
	}
	text = "hello world1"
	hash1, _ := vsha.Sha256Hash(text)
	if !VerifyPSS(pub, hash1, signature) {
		fmt.Println("签名不通过")
	}
}
