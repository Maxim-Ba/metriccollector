package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

type Signature struct {
	Key []byte
}

var signature Signature
var ErrKeyIsNotDefined = errors.New("key is not defined")

func New(key string) Signature {
	signature = Signature{
		Key: []byte(key),
	}
	return signature
}
func GetKey() string {
	return string(signature.Key)
}
func Get(src []byte) ([]byte, error) {
	if len(signature.Key) == 0 {
		return nil, ErrKeyIsNotDefined
	}
	h := hmac.New(sha256.New, signature.Key)
	_, err := h.Write(src)
	if err != nil {
		return nil, err
	}
	dst := h.Sum(nil)
	return dst, nil
}

func Check(dst []byte, bodySrc []byte) (error) {
	if len(signature.Key) == 0 {
		return ErrKeyIsNotDefined
	}
	h := hmac.New(sha256.New, signature.Key)
	_, err := h.Write(bodySrc)
	if err != nil {
		logger.LogError(err)
		return err
	}

	sign := h.Sum(nil)

	if !hmac.Equal(sign, dst) {
		return ErrInvalidSignature
	}

	return nil
}
