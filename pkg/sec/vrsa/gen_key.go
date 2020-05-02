package vrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// GenRsaKey gen_rsa_key
func GenRsaKey(bits int, priv io.Writer, pub io.Writer) (err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	err = pem.Encode(priv, block)
	if err != nil {
		return
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkiv, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkiv,
	}
	err = pem.Encode(pub, block)
	return
}
