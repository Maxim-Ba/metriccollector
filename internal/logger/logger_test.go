package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestSetLogLevel(t *testing.T) {
	// Создаем наблюдаемый логгер для проверки
	core, _ := observer.New(zapcore.InfoLevel)
	tempLogger := zap.New(core).Sugar()

	// Сохраняем оригинальный логгер и восстанавливаем после теста
	originalLogger := Sugar
	defer func() { Sugar = originalLogger }()

	Sugar = *tempLogger

	tests := []struct {
		name     string
		level    string
		expected zapcore.Level
	}{
		{"debug level", "debug", zap.DebugLevel},
		{"info level", "info", zap.InfoLevel},
		{"warn level", "warn", zap.WarnLevel},
		{"error level", "error", zap.ErrorLevel},
		{"unknown level", "unknown", zap.ErrorLevel}, // по умолчанию остается ErrorLevel
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.level)
			if atomicLevel.Level() != tt.expected {
				t.Errorf("SetLogLevel(%s) = %v, want %v", tt.level, atomicLevel.Level(), tt.expected)
			}
		})
	}
}

func TestSync(t *testing.T) {
	// Создаем наблюдаемый логгер для проверки ошибок
	core, recorded := observer.New(zapcore.ErrorLevel)
	tempLogger := zap.New(core).Sugar()

	// Сохраняем оригинальный логгер и восстанавливаем после теста
	originalLogger := Sugar
	defer func() { Sugar = originalLogger }()

	Sugar = *tempLogger

	// Тест 1: проверяем, что Sync не возвращает ошибку в нормальных условиях
	t.Run("successful sync", func(t *testing.T) {
		Sync()
		if recorded.Len() > 0 {
			t.Errorf("Sync() recorded unexpected error: %v", recorded.All())
		}
	})

}
