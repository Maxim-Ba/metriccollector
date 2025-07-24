package utils

import (
	"errors"
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	// Создаем ошибки один раз, чтобы можно было сравнивать указатели
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	err3 := errors.New("error3")

	testCases := []struct {
		name      string
		errorList []error
		err       error
		expected  bool
	}{
		{
			name:      "Error in list",
			errorList: []error{err1, err2},
			err:       err1, // Используем ту же ошибку
			expected:  true,
		},
		{
			name:      "Error not in list",
			errorList: []error{err1, err2},
			err:       err3, // Другая ошибка
			expected:  false,
		},
		{
			name:      "Empty list",
			errorList: []error{},
			err:       err1,
			expected:  false,
		},
		{
			name:      "Wrapped error",
			errorList: []error{err1},
			err:       errors.New("error1"), // Новая ошибка с тем же текстом
			expected:  false,                // Теперь ожидаем false, так как это разные экземпляры
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := contains(tc.errorList, tc.err)
			if result != tc.expected {
				t.Errorf("contains() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestRetryWrapperSuccess(t *testing.T) {
	counter := 0
	err := RetryWrapper(func() error {
		counter++
		return nil
	}, []error{errors.New("retryable error")})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if counter != 1 {
		t.Errorf("Expected 1 attempt, got %d", counter)
	}
}

func TestRetryWrapperRetryableError(t *testing.T) {
	retryableError := errors.New("retryable error")
	counter := 0
	startTime := time.Now()

	err := RetryWrapper(func() error {
		counter++
		if counter < 3 { // Первые две попытки возвращают ошибку
			return retryableError
		}
		return nil
	}, []error{retryableError})

	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}
	if counter != 3 {
		t.Errorf("Expected 3 attempts, got %d", counter)
	}

	elapsed := time.Since(startTime)
	if elapsed < 4*time.Second || elapsed > 5*time.Second {
		t.Errorf("Expected delays between retries, total time was %v", elapsed)
	}
}

func TestRetryWrapperNonRetryableError(t *testing.T) {
	retryableError := errors.New("retryable error")
	nonRetryableError := errors.New("non-retryable error")
	counter := 0

	err := RetryWrapper(func() error {
		counter++
		return nonRetryableError
	}, []error{retryableError})

	if !errors.Is(err, nonRetryableError) {
		t.Errorf("Expected nonRetryableError, got %v", err)
	}
	if counter != 1 {
		t.Errorf("Expected 1 attempt, got %d", counter)
	}
}

func TestRetryWrapperMaxRetries(t *testing.T) {
	retryableError := errors.New("retryable error")
	counter := 0

	err := RetryWrapper(func() error {
		counter++
		return retryableError
	}, []error{retryableError})

	if err == nil || err.Error() != "over max retry" {
		t.Errorf("Expected 'over max retry' error, got %v", err)
	}
	if counter != 3 {
		t.Errorf("Expected 3 attempts, got %d", counter)
	}
}
