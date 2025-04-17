package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

var saveInterval int
var localStoragePath string

func loadMetricsFromFile(path string) ([]*metrics.Metrics, error) {
	var metricsList []*metrics.Metrics

	// data, err := os.ReadFile(path)
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	// Проверка на пустые данные
	if len(data) == 0 {
		logger.LogInfo("Файл пуст, возвращаем пустой список метрик")
		return metricsList, nil
	}
	if err := json.Unmarshal(data, &metricsList); err != nil {
		logger.LogError(err)
		return nil, err
	}
	return metricsList, nil
}

func saveMetricsToFile(path string, metricsList *[]metrics.Metrics) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer file.Close()
	data, err := json.Marshal(metricsList)
	if err != nil {
		logger.LogError(err)
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.LogInfo("file saved")
	return err
}

func saveLoop() {
	if saveInterval == 0 {
		return
	}
	for {

		time.Sleep(time.Duration(saveInterval) * time.Second)
		paramsForGetAllMetrics := []*metrics.MetricDTOParams{}
		metricList, err := StorageInstance.GetMetrics(&paramsForGetAllMetrics)
		if err != nil {
			logger.LogError(err)
		}

		err = saveMetricsToFile(localStoragePath, metricList)
		if err != nil {
			logger.LogError(err)
		}
	}
}

func WithSyncLocalStorage (next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		if  saveInterval != 0 {
			next.ServeHTTP(res, r)
			return
		}
		next.ServeHTTP(res, r)
		logger.LogInfo("Metrics was saved")
		paramsForGetAllMetrics := []*metrics.MetricDTOParams{}
		metricList, err := StorageInstance.GetMetrics(&paramsForGetAllMetrics)
		if err != nil {
			logger.LogError(err)
		}

		err = saveMetricsToFile(localStoragePath, metricList)
		if err != nil {
			logger.LogError(err)
		}
	})
}
