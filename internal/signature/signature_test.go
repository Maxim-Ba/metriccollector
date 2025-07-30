package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

func TestNew(t *testing.T) {
	key := "test-key"
	sig := New(key, "")

	if string(sig.Key) != key {
		t.Errorf("Expected key %q, got %q", key, string(sig.Key))
	}

	// Проверяем, что глобальная переменная signature установлена
	if string(signature.Key) != key {
		t.Errorf("Expected global signature key %q, got %q", key, string(signature.Key))
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
			got := GetKey()
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
		got, err := Get(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !hmac.Equal(got, expectedHash) {
			t.Errorf("Expected hash %x, got %x", expectedHash, got)
		}
	})

	t.Run("Without key", func(t *testing.T) {
		New("", "") // Сбрасываем ключ
		_, err := Get(data)
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
			err := Check(tt.signature, tt.data)

			if err != tt.expectedErr {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

// Вспомогательная функция для вычисления HMAC
func computeHMAC(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return h.Sum(nil)
}
