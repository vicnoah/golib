package vaes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	text := "123"                                        // 你要加密的数据
	AesKey := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR") // 对称秘钥长度必须是16的倍数

	fmt.Printf("明文: %s\n秘钥: %s\n", text, string(AesKey))
	encrypted, err := CBCEncrypt([]byte(text), AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))
	//encrypteds, _ := base64.StdEncoding.DecodeString("j4H4Tv5VcXv0oNLwB/fr+g==")
	origin, err := CBCDecrypt(encrypted, AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("解密后明文: %s\n", string(origin))
}
