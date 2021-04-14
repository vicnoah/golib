package rand

import (
	"crypto/rand"
)

var strstr = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GetString 生成随机字符串
func GetString(size int) string {
	data := make([]byte, size)
	out := make([]byte, size)
	buffer := len(strstr)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	for id, key := range data {
		x := byte(int(key) % buffer)
		out[id] = strstr[x]
	}
	return string(out)
}
