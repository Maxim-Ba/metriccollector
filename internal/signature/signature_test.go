package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestNew(t *testing.T) {
	key := "test-key"
	sig := New(key, "")

	if string(sig.Key) != key {
		t.Errorf("Expected key %q, got %q", key, string(sig.Key))
	}

	if string(Instance.Key) != key {
		t.Errorf("Expected global signature key %q, got %q", key, string(Instance.Key))
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"Empty key", "", ""},
		{"Non-empty key", "test-key", "test-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.key, "")
			got := Instance.GetKey()
			if got != tt.expected {
				t.Errorf("Expected key %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Тестовые данные
	key := "test-key"
	data := []byte("test-data")
	expectedHash := computeHMAC(data, key)

	t.Run("With key", func(t *testing.T) {
		New(key, "")
		got, err := Instance.Get(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !hmac.Equal(got, expectedHash) {
			t.Errorf("Expected hash %x, got %x", expectedHash, got)
		}
	})

	t.Run("Without key", func(t *testing.T) {
		New("", "") // Сбрасываем ключ
		_, err := Instance.Get(data)
		if err != ErrKeyIsNotDefined {
			t.Errorf("Expected error %v, got %v", ErrKeyIsNotDefined, err)
		}
	})
}

func TestCheck(t *testing.T) {
	key := "test-key"
	data := []byte("test-data")
	validHash := computeHMAC(data, key)
	invalidHash := computeHMAC(data, "wrong-key")

	tests := []struct {
		name        string
		key         string
		signature   []byte
		data        []byte
		expectedErr error
	}{
		{"Valid signature", key, validHash, data, nil},
		{"Invalid signature", key, invalidHash, data, ErrInvalidSignature},
		{"Empty key", "", validHash, data, ErrKeyIsNotDefined},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.key, "")
			err := Instance.Check(tt.signature, tt.data)

			if err != tt.expectedErr {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func computeHMAC(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return h.Sum(nil)
}

func TestEncrypt(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	testData := []byte("test data for encryption")

	t.Run("Successful encryption", func(t *testing.T) {
		Instance = &Signature{PublicKey: publicKey}

		encrypted, err := Instance.Encrypt(testData)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		if string(encrypted) == string(testData) {
			t.Error("Encrypted data should not match original data")
		}

		if len(encrypted) == 0 {
			t.Error("Encrypted data should not be empty")
		}
	})

	t.Run("Encrypt without public key", func(t *testing.T) {
		Instance = &Signature{PublicKey: nil}

		_, err := Instance.Encrypt(testData)
		if err != ErrKeyIsNotDefined {
			t.Errorf("Expected error %v, got %v", ErrKeyIsNotDefined, err)
		}
	})

	t.Run("Encrypt empty data", func(t *testing.T) {
		Instance = &Signature{PublicKey: publicKey}

		_, err := Instance.Encrypt([]byte{})
		if err != nil {
			t.Errorf("Encrypting empty data should not return error, got %v", err)
		}
	})
}

func TestDecrypt(t *testing.T) {
	// Генерируем тестовую пару RSA ключей
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Тестовые данные
	testData := []byte("test data for decryption")

	t.Run("Successful decryption", func(t *testing.T) {
		Instance = &Signature{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}

		encrypted, err := Instance.Encrypt(testData)
		if err != nil {
			t.Fatalf("Setup failed: couldn't encrypt test data: %v", err)
		}

		decrypted, err := Instance.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt failed: %v", err)
		}

		if string(decrypted) != string(testData) {
			t.Errorf("Decrypted data doesn't match original. Got %q, want %q", decrypted, testData)
		}
	})

	t.Run("Decrypt without private key", func(t *testing.T) {
		Instance = &Signature{PublicKey: publicKey}
		encrypted, err := Instance.Encrypt(testData)
		if err != nil {
			t.Fatalf("Setup failed: couldn't encrypt test data: %v", err)
		}

		Instance = &Signature{PrivateKey: nil}
		_, err = Instance.Decrypt(encrypted)
		if err != ErrKeyIsNotDefined {
			t.Errorf("Expected error %v, got %v", ErrKeyIsNotDefined, err)
		}
	})

	t.Run("Decrypt corrupted data", func(t *testing.T) {
		Instance = &Signature{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}

		encrypted, err := Instance.Encrypt(testData)
		if err != nil {
			t.Fatalf("Setup failed: couldn't encrypt test data: %v", err)
		}

		// Повреждаем зашифрованные данные
		if len(encrypted) > 10 {
			encrypted[5] ^= 0xFF // Инвертируем байт
		}

		_, err = Instance.Decrypt(encrypted)
		if err == nil {
			t.Error("Expected error when decrypting corrupted data, got nil")
		}
	})

	t.Run("Decrypt empty data", func(t *testing.T) {
		Instance = &Signature{PrivateKey: privateKey}

		_, err := Instance.Decrypt([]byte{})
		if err != nil {
			t.Errorf("Decrypting empty data should not return error, got %v", err)
		}
	})
}

func TestGetPubKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	tests := []struct {
		name     string
		setup    func()
		expected *rsa.PublicKey
	}{
		{
			name: "Public key exists",
			setup: func() {
				Instance = &Signature{PublicKey: publicKey}
			},
			expected: publicKey,
		},
		{
			name: "Public key is nil",
			setup: func() {
				Instance = &Signature{PublicKey: nil}
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := Instance.GetPubKey()

			if got != tt.expected {
				t.Errorf("Expected public key %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestGetPrivKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	tests := []struct {
		name     string
		setup    func()
		expected *rsa.PrivateKey
	}{
		{
			name: "Private key exists",
			setup: func() {
				Instance = &Signature{PrivateKey: privateKey}
			},
			expected: privateKey,
		},
		{
			name: "Private key is nil",
			setup: func() {
				Instance = &Signature{PrivateKey: nil}
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := Instance.GetPrivKey()

			if got != tt.expected {
				t.Errorf("Expected private key %v, got %v", tt.expected, got)
			}
		})
	}
}
