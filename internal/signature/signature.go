package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

type Signature struct {
	Key        []byte
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var (
	signature          Signature
	ErrKeyIsNotDefined = errors.New("key is not defined")
	ErrInvalidKey      = errors.New("invalid key format")
)

func New(key string, cryptoKeyPath string) Signature {
	if cryptoKeyPath != "" {
		b, err := os.ReadFile(cryptoKeyPath)
		if err != nil {
			logger.LogError(err)
			return signature
		}

		block, _ := pem.Decode(b)
		if block == nil {
			logger.LogError(ErrInvalidKey)
			return signature
		}

		// для сервера пытаемся загрузить приватный ключ
		if privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			signature = Signature{
				PrivateKey: privKey,
				PublicKey:  &privKey.PublicKey,
			}
			return signature
		}

		// для агента пытаемся загрузить публичный ключ
		if pubKey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if rsaPubKey, ok := pubKey.(*rsa.PublicKey); ok {
				signature = Signature{
					PublicKey: rsaPubKey,
				}
				return signature
			}
		}

		logger.LogError(ErrInvalidKey)
		return signature
	}
	signature = Signature{
		Key: []byte(key),
	}
	return signature
}

func GetKey() string {
	return string(signature.Key)
}
func GetPubKey() *rsa.PublicKey {
	return signature.PublicKey
}
func GetPrivKey() *rsa.PrivateKey {
	return signature.PrivateKey
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

func Check(dst []byte, bodySrc []byte) error {
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

// -------- assym

func Encrypt(data []byte) ([]byte, error) {
	if signature.PublicKey == nil {
		return nil, ErrKeyIsNotDefined
	}

	hash := sha256.New()
	msgLen := len(data)
	step := signature.PublicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlock, err := rsa.EncryptOAEP(
			hash,
			rand.Reader,
			signature.PublicKey,
			data[start:finish],
			nil,
		)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlock...)
	}

	return encryptedBytes, nil
}

func Decrypt(data []byte) ([]byte, error) {
	if signature.PrivateKey == nil {
		return nil, ErrKeyIsNotDefined
	}

	hash := sha256.New()
	keySize := signature.PrivateKey.Size()
	var decryptedBytes []byte

	for start := 0; start < len(data); start += keySize {
		finish := start + keySize
		if finish > len(data) {
			finish = len(data)
		}

		decryptedBlock, err := rsa.DecryptOAEP(
			hash,
			rand.Reader,
			signature.PrivateKey,
			data[start:finish],
			nil,
		)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlock...)
	}

	return decryptedBytes, nil
}
