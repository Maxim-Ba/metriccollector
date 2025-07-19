package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage/database/postgres"
)

var saveInterval int
var localStoragePath string
var databaseDSN string

func loadMetricsFromFile(path string) ([]*metrics.Metrics, error) {
	var metricsList []*metrics.Metrics

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.LogError(err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	// Проверка на пустые данные
	if len(data) == 0 {
		logger.LogInfo("file is empty")
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
	defer func() {
		if err = file.Close(); err != nil {
			logger.LogError(err)
		}
	}()
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
		if databaseDSN != "" {
			err = postgres.SaveMetricsToDB(metricList, db)
		} else {
			err = saveMetricsToFile(localStoragePath, metricList)
		}
		if err != nil {
			logger.LogError(err)
		}
	}
}

func WithSyncLocalStorage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		if saveInterval != 0 {
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
