package signature

import (
	"crypto/hmac"
	"crypto/sha256"
)

var hashKey = []byte("")

func SetKey(key string) {
	hashKey = []byte(key)
}
func GetKey() string {
	return string(hashKey)
}
func Get(src []byte) ([]byte, error) {
	h := hmac.New(sha256.New, hashKey)
	_, err := h.Write(src)
	if err != nil {
		return nil, err
	}
	dst := h.Sum(nil)
	return dst, nil
}

func Check(src []byte) bool {
	h := hmac.New(sha256.New, hashKey)
	sign := h.Sum(nil)
	_, err := h.Write(src)
	if err != nil {
		return false
	}
	return hmac.Equal(sign, src)
}
