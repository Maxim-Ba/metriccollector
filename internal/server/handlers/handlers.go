package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	metricsService "github.com/Maxim-Ba/metriccollector/internal/server/services/metric"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)

func GetAllHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getAllHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}


	html, err := metricsService.GetAll(storage.StorageInstance)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(""))
		return
	}
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(html))
}
func GetOneHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getOneHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogInfo("B------------")
		logger.LogInfo(err)
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(""))
			return
	}
	requestMetric, err := parseMetric(&buf)

	if err != nil {

		logger.LogInfo("C------------")
		logger.LogInfo(err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(""))
		return
	}

	metricsNames :=[]string{requestMetric.ID}

	responseMetrics, err := metricsService.Get(storage.StorageInstance, &metricsNames)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(""))
		return
	}

		body, err := json.Marshal(responseMetrics)
	
	if err != nil {
		logger.LogInfo("D------------")
		logger.LogInfo(err)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(""))
		return
}
	res.Header().Set("Content-Type", " application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func UpdateHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("updateHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogInfo("1------------")
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogInfo("2------------")
		logger.LogInfo(err)

			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(""))
			return
	}
	metric, err := parseMetric(&buf)

	if err != nil {
		logger.LogInfo("3------------")
		logger.LogInfo(err)
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		res.Write([]byte(""))
		return
	}
	logger.LogInfo("------metric------")
	logger.LogMetric(metric)

	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		logger.LogInfo("4------------")
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(""))
}

func checkForAllowedMethod(req *http.Request, allowedMethod []string) error {
	if !(slices.Contains(allowedMethod, req.Method)) {
		return fmt.Errorf("not allowed method")
	}
	return nil
}
func parseMetric(buf *bytes.Buffer) (metrics.Metrics, error) {
	var metric metrics.Metrics
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		logger.LogInfo(err)
		return metrics.Metrics{}, ErrNoMetricName
	}

	if metric.MType != "gauge" && metric.MType != "counter" {
		return metrics.Metrics{}, ErrNoMetricsType
	}
	if metric.ID == "" {
		return metrics.Metrics{}, ErrNoMetricName
	}
	return metric, nil
}
func parseMetrics(buf *bytes.Buffer) ([]metrics.Metrics, error) {
	var metricsList []metrics.Metrics
	if err := json.Unmarshal(buf.Bytes(), &metricsList); err != nil {
		logger.LogInfo(err)
		return nil, ErrNoMetricName
	}

	for _, metric := range metricsList {
		if metric.MType != "gauge" && metric.MType != "counter" {
			return nil, ErrNoMetricsType
		}
		if metric.ID == "" {
			return nil, ErrNoMetricName
		}
	}

	return metricsList, nil
}
