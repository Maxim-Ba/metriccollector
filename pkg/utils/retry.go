package utils

import (
	"errors"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

// RetryWrapper вызывает функцию, повторяя попытку в случае ошибки из переданного списка.
func RetryWrapper(action func() error, retries int, errorList []error) error {
    for i := 0; i < retries; i++ {
        err := action()
        if err != nil && contains(errorList, err) {
            logger.LogError("retry %d/%d: %v\n", i+1, retries, err)
            time.Sleep(time.Duration(2*i+1) * time.Second) // ждем, чтобы дать время для подготовки
            continue
        }
        return err
    }
    return errors.New("over max retry")
}

// contains проверяет наличие ошибки в списке
func contains(errorList []error, err error) bool {
    for _, e := range errorList {
        if errors.Is(e, err) {
            return true
        }
    }
    return false
}
