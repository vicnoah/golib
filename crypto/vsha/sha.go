package vsha

import "crypto/sha256"

// Sha256Hash 生成sha256 hash
func Sha256Hash(text string) (digest string, err error) {
	msg := []byte(text)

	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return
	}
	digest = string(msgHash.Sum(nil))
	return
}
