package utils

import (
	"errors"
	"fmt"
	"time"
)

// RetryWrapper вызывает функцию, повторяя попытку в случае ошибки из переданного списка.
func RetryWrapper(action func() error, retries int, errorList []error) error {
    for i := 0; i < retries; i++ {
        err := action()
        if err != nil && contains(errorList, err) {
            fmt.Printf("Возникла ошибка из списка, повтор %d/%d: %v\n", i+1, retries, err)
            time.Sleep(time.Duration(2*i+1) * time.Second) // ждем, чтобы дать время для подготовки
            continue
        }
        return err
    }
    return errors.New("превышено количество повторов")
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
