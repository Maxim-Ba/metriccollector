package signature

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

var hashKey = []byte("")

func SetKey(key string) {
	logger.LogInfo("key" , key, "-1-1-1-1-1-1-")

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

func Check(dst []byte, bodySrc []byte) bool {
	h := hmac.New(sha256.New, hashKey)
	_, err := h.Write(bodySrc)
	if err != nil {
		logger.LogError(err)
		return false
	}

	sign := h.Sum(nil)
	logger.LogInfo(sign, dst)
	return hmac.Equal(sign, dst)
}
