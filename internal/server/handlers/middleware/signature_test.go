package middleware

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/signature"
	"github.com/stretchr/testify/require"
)

func TestSignatureHandle(t *testing.T) {
	originalKey := signature.GetKey()
	defer func() {
		signature.New(originalKey)
	}()

	signature.New("test-key")

	// Генерируем валидную подпись для тестовых данных
	validData := []byte("test data")
	validHash, err := signature.Get(validData)
	if err != nil {
		t.Fatalf("Failed to generate valid hash: %v", err)
	}
	validHashBase64 := base64.StdEncoding.EncodeToString(validHash)

	tests := []struct {
		name           string
		requestHeader  map[string]string
		requestBody    string
		expectStatus   int
		expectResponse string
		checkResponse  bool
		setup          func()
		cleanup        func()
	}{
		{
			name:          "No signature header - pass through",
			requestHeader: nil,
			requestBody:   "test data",
			expectStatus:  http.StatusOK,
		},
		{
			name: "Valid signature",
			requestHeader: map[string]string{
				"HashSHA256": validHashBase64,
			},
			requestBody:    "test data",
			expectStatus:   http.StatusOK,
			checkResponse:  true,
			expectResponse: "test response",
		},
		{
			name: "Invalid base64 signature",
			requestHeader: map[string]string{
				"HashSHA256": "invalid base64!!!",
			},
			requestBody:    "test data",
			expectStatus:   http.StatusBadRequest,
			expectResponse: "invalid base64 encoding\n",
		},
		{
			name: "Invalid signature",
			requestHeader: map[string]string{
				"HashSHA256": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", // неверная подпись
			},
			requestBody:    "test data",
			expectStatus:   http.StatusBadRequest,
			expectResponse: "\n",
		},
		{
			name: "No key configured - pass through",
			requestHeader: map[string]string{
				"HashSHA256": validHashBase64,
			},
			requestBody:    "test data",
			expectStatus:   http.StatusOK,
			expectResponse: "test response",
			setup: func() {
				signature.New("") // очищаем ключ
			},
			cleanup: func() {
				signature.New("test-key") // восстанавливаем ключ
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkResponse {
					var body []byte
					body, err = io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}
					if string(body) != tt.requestBody {
						t.Errorf("Expected request body %q, got %q", tt.requestBody, string(body))
					}
				}
				_, err = w.Write([]byte("test response"))
				if err != nil {
					require.NoError(t, err)
				}
			})

			wrappedHandler := SignatureHandle(handler)

			req := httptest.NewRequest("POST", "http://example.com", bytes.NewBufferString(tt.requestBody))
			if tt.requestHeader != nil {
				for k, v := range tt.requestHeader {
					req.Header.Set(k, v)
				}
			}

			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectStatus)
			}

			if tt.expectResponse != "" && rr.Body.String() != tt.expectResponse {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectResponse)
			}

			if tt.checkResponse {
				if rr.Header().Get("HashSHA256") == "" {
					t.Error("expected HashSHA256 header in response, got none")
				}
			}
		})
	}
}

func TestHashResponseWriter(t *testing.T) {
	originalKey := signature.GetKey()
	defer func() {
		signature.New(originalKey)
	}()
	signature.New("test-key")

	t.Run("With key set - should add signature header", func(t *testing.T) {
		rr := httptest.NewRecorder()

		hw := &hashResponseWriter{
			ResponseWriter: rr,
		}

		// Тестовые данные
		testData := []byte("test data")
		expectedHash, err := signature.Get(testData)
		if err != nil {
			t.Fatalf("Failed to generate expected hash: %v", err)
		}
		expectedHashBase64 := base64.StdEncoding.EncodeToString(expectedHash)

		// Записываем данные
		size, err := hw.Write(testData)
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		// Проверяем размер записанных данных
		if size != len(testData) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(testData), size)
		}

		// Проверяем записанные данные
		if rr.Body.String() != string(testData) {
			t.Errorf("Expected body %q, got %q", string(testData), rr.Body.String())
		}

		// Проверяем заголовок подписи
		actualHash := hw.Header().Get("HashSHA256")
		if actualHash != expectedHashBase64 {
			t.Errorf("Expected hash header %q, got %q", expectedHashBase64, actualHash)
		}
	})

}
