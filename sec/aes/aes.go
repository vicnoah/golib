// Package aes aes encryption and decryption
// Because the Go standard library has incomplete support for AES-192, AES-256, only aes-cbc-128 is allowed
// blockSize 16 = aes-cbc-128, 24 = aes-cbc-192, 32 = aes-cbc-256
// Related Links https://studygolang.com/articles/7302
// https://blog.csdn.net/xiaohu50/article/details/51682849
// https://www.jianshu.com/p/3741458695d2
package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

const (
	// 加解密时用于补位的字符
	_PAD = byte('0')
	// 密码长度
	_BLOCK_SIZE = 16
)


var (
	// PKCS7 errors.
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// CBCDecrypt Decrypt aes data
//接收一个16进制编码的字节切片
//参数1:16进制aes加密数据,参数2:aes加密key,参数3:aes加密偏移向量
func CBCDecrypt(hexStr []byte, key []byte, iv []byte) (data []byte, err error) {
	cipherText, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return
	}
	block, err := aes.NewCipher(PaddingLeft(key, _PAD, _BLOCK_SIZE)) //选择加密算法
	if err != nil {
		return
	}
	//blockModel := cipher.NewCBCDecrypter(block, keyBytes)
	blockModel := cipher.NewCBCDecrypter(block, PaddingLeft(iv, _PAD, _BLOCK_SIZE))
	plantText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plantText, cipherText)
	//plantText = PKCS7UnPadding(plantText, block.BlockSize())
	data, err = pkcs7Unpad(plantText, block.BlockSize())
	return
}

// CBCEncrypt Returns a hexadecimal encoded aes encrypted string
func CBCEncrypt(plantText []byte, key []byte, iv []byte) (cryptBytes []byte, err error) {
	block, err := aes.NewCipher(PaddingLeft(key, _PAD, _BLOCK_SIZE)) //选择加密算法
	if err != nil {
		return
	}
	//plantText = PKCS7Padding(plantText, block.BlockSize())
	plantText, err = pkcs7Pad(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, PaddingLeft(iv, _PAD, _BLOCK_SIZE))
	cipherText := make([]byte, len(plantText))
	blockModel.CryptBlocks(cipherText, plantText)
	cryptBytes = []byte(hex.EncodeToString(cipherText))
	return
}

// PKCS7Padding aes padding algorithm
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS7UnPadding aes unpadding算法
func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// PaddingLeft Password left complement algorithm
func PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

// PKCS7 padding.

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blockSize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}