package client

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient_SendMetrics(t *testing.T) {
	type args struct {
		metrics []*metrics.Metrics
	}
	tests := []struct {
		name        string
		client      *HTTPClient
		args        args
		wantErr     bool
		wantHeaders map[string]string
	}{
		{
			name: "Valid Content-Type",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{
					{},
				},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Valid Content-Encoding",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{{}},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Content-Encoding": "gzip",
			},
		},
		{
			name: "Valid Accept-Encoding",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{{}},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Accept-Encoding": "gzip",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInstance := signature.Instance
			defer func() {
				signature.Instance = originalInstance
			}()

			// Инициализируем новый Instance для теста
			signature.New("test-key", "")

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, value := range tt.wantHeaders {
					if r.Header.Get(key) != value {
						t.Errorf("Expected header %s to be %s, got %s", key, value, r.Header.Get(key))
					}
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			address = server.URL[7:] // Remove "http://"
			if err := tt.client.SendMetrics(tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.SendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	testAddress := "localhost:8080"
	client := NewClient(testAddress)

	if client == nil {
		t.Fatal("Expected non-nil HTTPClient, got nil")
	}

	if client.httpClient == nil {
		t.Fatal("Expected non-nil http.Client, got nil")
	}

	if address != testAddress {
		t.Errorf("Expected address to be %s, got %s", testAddress, address)
	}

	expectedTimeout := 10 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}

func TestSendMetricsWithBatch(t *testing.T) {
	originalInstance := signature.Instance
	defer func() {
		signature.Instance = originalInstance
	}()

	// Инициализируем новый Instance для теста
	signature.New("test-key", "")

	// Создаем тестовые данные
	testMetrics := []*metrics.Metrics{
		{ID: "test1", MType: "gauge", Value: new(float64)},
		{ID: "test2", MType: "counter", Delta: new(int64)},
	}
	*testMetrics[0].Value = 123.45
	*testMetrics[1].Delta = 42

	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем путь
		assert.Equal(t, "/updates/", r.URL.Path)

		// Проверяем заголовки
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Декомпрессируем тело
		gz, err := gzip.NewReader(r.Body)
		require.NoError(t, err)
		defer gz.Close()

		var receivedMetrics []*metrics.Metrics
		err = json.NewDecoder(gz).Decode(&receivedMetrics)
		require.NoError(t, err)

		// Проверяем полученные метрики
		assert.Equal(t, testMetrics[0].ID, receivedMetrics[0].ID)
		assert.Equal(t, *testMetrics[0].Value, *receivedMetrics[0].Value)
		assert.Equal(t, testMetrics[1].ID, receivedMetrics[1].ID)
		assert.Equal(t, *testMetrics[1].Delta, *receivedMetrics[1].Delta)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Создаем клиент с адресом тестового сервера
	client := NewClient(ts.URL[7:]) // ts.URL[7:] чтобы убрать "http://"

	// Вызываем тестируемый метод
	err := client.SendMetricsWithBatch(testMetrics)
	assert.NoError(t, err)
}

func TestSendMetricsWithBatch_ErrorCases(t *testing.T) {
	originalInstance := signature.Instance
	defer func() {
		signature.Instance = originalInstance
	}()

	// Инициализируем новый Instance для теста
	signature.New("test-key", "")
	tests := []struct {
		name           string
		responseStatus int
		expectedError  error
	}{
		{
			name:           "Internal Server Error",
			responseStatus: http.StatusInternalServerError,
			expectedError:  ErrServerInternalError,
		},
		{
			name:           "Request Timeout",
			responseStatus: http.StatusRequestTimeout,
			expectedError:  ErrRequestTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый сервер, который возвращает нужный статус
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
			}))
			defer ts.Close()

			// Создаем клиент
			client := NewClient(ts.URL[7:])

			// Вызываем метод с пустыми метриками (нам важна только ошибка)
			err := client.SendMetricsWithBatch([]*metrics.Metrics{})
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
