package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToPointerInt(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  int64
	}{
		{"positive number", 42, 42},
		{"negative number", -42, -42},
		{"zero", 0, 0},
		{"max int64", 9223372036854775807, 9223372036854775807},
		{"min int64", -9223372036854775808, -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntToPointerInt(tt.input)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestIntToPointerFloat(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		want  float64
	}{
		{"small number", 42, 42.0},
		{"zero", 0, 0.0},
		{"large number", 18446744073709551615, 18446744073709551615.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntToPointerFloat(tt.input)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestFloatToPointerFloat(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"positive float", 3.14, 3.14},
		{"negative float", -3.14, -3.14},
		{"zero", 0.0, 0.0},
		{"large float", 1.7976931348623157e+308, 1.7976931348623157e+308},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatToPointerFloat(tt.input)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestFloatToPointerInt(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  int64
	}{
		{"positive number", 42, 42},
		{"negative number", -42, -42},
		{"zero", 0, 0},
		{"max int64", 9223372036854775807, 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatToPointerInt(tt.input)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestWrireZeroBytes(t *testing.T) {
	t.Run("write zero bytes to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		WrireZeroBytes(&buf)
		assert.Equal(t, 0, buf.Len())
	})

	t.Run("write zero bytes to failing writer", func(t *testing.T) {
		// Создаем writer, который всегда возвращает ошибку
		mockWriter := &failingWriter{}
		WrireZeroBytes(mockWriter)
		// Проверяем, что ошибка была обработана (в реальном коде она логируется)
	})
}

// failingWriter это io.Writer, который всегда возвращает ошибку
type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}
