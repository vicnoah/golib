// Package vmd5 md5扩展库
package vmd5

import (
	"crypto/md5"
	"encoding/hex"
)

// Sum 32位md5值计算
func Sum(data []byte) ([]byte, error) {
	//生成32位md5
	myMd5 := md5.New()
	_, err := myMd5.Write(data)
	if err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(myMd5.Sum(nil))), nil
}
