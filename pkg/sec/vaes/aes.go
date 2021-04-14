// Package vaes aes encryption and decryption
// Because the Go standard library has incomplete support for AES-192, AES-256, only aes-cbc-128 is allowed
// blockSize 16 = aes-cbc-128, 24 = aes-cbc-192, 32 = aes-cbc-256
// Related Links https://studygolang.com/articles/7302
// https://blog.csdn.net/xiaohu50/article/details/51682849
// https://www.jianshu.com/p/3741458695d2
package vaes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

const (
	// 加解密时用于补位的字符
	_PAD = byte('0')
	// 密码长度
	_BLOCK_SIZE = 16
)

// CBCEncrypt Returns a hexadecimal encoded aes encrypted string
func CBCEncrypt(plantText []byte, key []byte, iv []byte) (cryptBytes []byte, err error) {
	block, err := aes.NewCipher(PaddingLeft(key, _PAD, _BLOCK_SIZE)) //选择加密算法
	if err != nil {
		return
	}
	plantText = Padding(PKCS7_PADDING, plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, PaddingLeft(iv, _PAD, _BLOCK_SIZE))
	cipherText := make([]byte, len(plantText))
	blockModel.CryptBlocks(cipherText, plantText)
	cryptBytes = []byte(hex.EncodeToString(cipherText))
	return
}

// CBCDecrypt Decrypt aes data
// 接收一个16进制编码的字节切片
// 参数1:16进制aes加密数据,参数2:aes加密key,参数3:aes加密偏移向量
func CBCDecrypt(hexStr []byte, key []byte, iv []byte) (data []byte, err error) {
	cipherText, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return
	}
	block, err := aes.NewCipher(PaddingLeft(key, _PAD, _BLOCK_SIZE)) //选择加密算法
	if err != nil {
		return
	}
	blockModel := cipher.NewCBCDecrypter(block, PaddingLeft(iv, _PAD, _BLOCK_SIZE))
	plantText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plantText, cipherText)
	data = UnPadding(PKCS7_PADDING, plantText)
	return
}
