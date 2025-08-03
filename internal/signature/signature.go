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
	Instance           *Signature
	ErrKeyIsNotDefined = errors.New("key is not defined")
	ErrInvalidKey      = errors.New("invalid key format")
)

func New(key string, cryptoKeyPath string) *Signature {
	sig := &Signature{}
	Instance = sig
	if cryptoKeyPath != "" {
		b, err := os.ReadFile(cryptoKeyPath)
		if err != nil {
			logger.LogError(err)
			return sig
		}

		block, _ := pem.Decode(b)
		if block == nil {
			logger.LogError(ErrInvalidKey)
			return sig
		}

		// для сервера пытаемся загрузить приватный ключ
		if privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			// signature = Signature{
			sig.PrivateKey = privKey
			sig.PublicKey = &privKey.PublicKey
			// }
			return sig
		}

		// для агента пытаемся загрузить публичный ключ
		if pubKey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if rsaPubKey, ok := pubKey.(*rsa.PublicKey); ok {
				// signature = Signature{
				sig.PublicKey = rsaPubKey
				// }
				return sig
			}
		}

		logger.LogError(ErrInvalidKey)
		return sig
	}
	// signature = Signature{
	sig.Key = []byte(key)
	// }
	return sig
}

func (s *Signature) GetKey() string {
	return string(s.Key)
}

func (s *Signature) GetPubKey() *rsa.PublicKey {
	return s.PublicKey
}

func (s *Signature) GetPrivKey() *rsa.PrivateKey {
	return s.PrivateKey
}

func (s *Signature) Get(src []byte) ([]byte, error) {
	if len(s.Key) == 0 {
		return nil, ErrKeyIsNotDefined
	}
	h := hmac.New(sha256.New, s.Key)
	_, err := h.Write(src)
	if err != nil {
		return nil, err
	}
	dst := h.Sum(nil)
	return dst, nil
}

func (s *Signature) Check(dst []byte, bodySrc []byte) error {
	if len(s.Key) == 0 {
		return ErrKeyIsNotDefined
	}
	h := hmac.New(sha256.New, s.Key)
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

func (s *Signature) Encrypt(data []byte) ([]byte, error) {
	if s.PublicKey == nil {
		return nil, ErrKeyIsNotDefined
	}

	hash := sha256.New()
	msgLen := len(data)
	step := s.PublicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlock, err := rsa.EncryptOAEP(
			hash,
			rand.Reader,
			s.PublicKey,
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

func (s *Signature) Decrypt(data []byte) ([]byte, error) {
	if s.PrivateKey == nil {
		return nil, ErrKeyIsNotDefined
	}

	hash := sha256.New()
	keySize := s.PrivateKey.Size()
	var decryptedBytes []byte

	for start := 0; start < len(data); start += keySize {
		finish := start + keySize
		if finish > len(data) {
			finish = len(data)
		}

		decryptedBlock, err := rsa.DecryptOAEP(
			hash,
			rand.Reader,
			s.PrivateKey,
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
