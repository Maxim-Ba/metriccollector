package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

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
func  UpdateHandlerByURLParams (res http.ResponseWriter, req *http.Request) {
	fmt.Print("UpdateHandlerByURLParams \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	urlString := req.URL.Path //  /update/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/update/")
	parameters := strings.Split(params, "/")
	metric, err := metricRecord(parameters)
	logger.LogInfo("urlString------------")
	logger.LogInfo(urlString)

	if err != nil {
		logger.LogInfo("UpdateHandlerByURLParams3------------")
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
	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		logger.LogInfo("UpdateHandlerByURLParams4------------")
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(params))
	res.Write([]byte(""))
}
func metricRecord(parameters []string) (metrics.Metrics, error) {
	if len(parameters) != 3 {
		return metrics.Metrics{}, ErrNoMetricName
	}
		if parameters[0] != "gauge" && parameters[0] != "counter" {
			return metrics.Metrics{}, ErrNoMetricsType
		
		}
		if parameters[1] == "" {
			return metrics.Metrics{}, ErrNoMetricName
		}
		metricType := parameters[0]
		metricName := parameters[1]
		var value float64
		var delta int64
		var err error
		if parameters[0] == "gauge" {
			value, err = strconv.ParseFloat(parameters[2], 64)
			} else {
				delta, err = strconv.ParseInt(parameters[2], 10,64)

			}
			if err != nil {
				return metrics.Metrics{}, ErrWrongValue
			}
		
		return metrics.Metrics{
			MType: metricType,
			ID: metricName,
			Value:      &value,
			Delta: &delta,
		}, nil
	
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
