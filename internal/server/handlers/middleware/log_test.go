package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingResponseWriter_Write(t *testing.T) {
	mockRW := httptest.NewRecorder()
	responseData := &responseData{}

	lw := loggingResponseWriter{
		ResponseWriter: mockRW,
		responseData:   responseData,
	}

	data := []byte("test data")
	size, err := lw.Write(data)

	require.NoError(t, err)
	assert.Equal(t, len(data), size)
	assert.Equal(t, len(data), responseData.size)
	assert.Equal(t, data, mockRW.Body.Bytes())
}

func TestLoggingResponseWriter_WriteHeader(t *testing.T) {
	mockRW := httptest.NewRecorder()
	responseData := &responseData{}

	lw := loggingResponseWriter{
		ResponseWriter: mockRW,
		responseData:   responseData,
	}

	// Устанавливаем статус код
	statusCode := http.StatusNotFound
	lw.WriteHeader(statusCode)

	assert.Equal(t, statusCode, responseData.status)
	assert.Equal(t, statusCode, mockRW.Code)
}

func TestWithLogging_MultipleWrites(t *testing.T) {
	// Создаем наблюдаемый логгер (хотя в этом тесте мы не проверяем логи)
	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	logger.Sugar = *observedLogger.Sugar()

	// Тестовый обработчик с несколькими записями
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("part1 "))
		if err != nil {
			fmt.Print(err)
		}
		_, err = w.Write([]byte("part2"))
		if err != nil {
			fmt.Print(err)
		}
	})

	// Оборачиваем обработчик в middleware
	wrappedHandler := WithLogging(testHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "http://example.com/multi", nil)
	rr := httptest.NewRecorder()

	// Выполняем запрос
	wrappedHandler.ServeHTTP(rr, req)

	// Проверяем что размер ответа корректно суммировался
	assert.Equal(t, "part1 part2", rr.Body.String())
	assert.Equal(t, 11, rr.Body.Len())
}
