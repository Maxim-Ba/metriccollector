package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	metricsService "github.com/Maxim-Ba/metriccollector/internal/server/services/metric"
	storageService "github.com/Maxim-Ba/metriccollector/internal/server/services/starage"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

// GetAllHandler handles HTTP GET requests to retrieve all metrics.
// Returns an HTML page listing all metrics in storage.
// Responds with appropriate HTTP status codes for errors.
func GetAllHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("getAllHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)

		utils.WrireZeroBytes(res)
		return
	}

	html, err := metricsService.GetAll(storage.StorageInstance)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)

		utils.WrireZeroBytes(res)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	_, err = res.Write([]byte(html))

	if err != nil {
		logger.LogError(err)
		return
	}
}

// GetOneHandlerByParams handles HTTP GET requests to retrieve a single metric via URL parameters.
// Expected URL format: /value/<type>/<name>.
// Returns the metric value as plain text.
// Responds with HTTP 404 if metric is not found, or 400 for bad requests.
func GetOneHandlerByParams(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("GetOneHandlerByParams")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)

		utils.WrireZeroBytes(res)
		return
	}
	urlString := req.URL.Path //  /value/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/value/")
	parameters := strings.Split(params, "/")
	name := parameters[1]
	metricParams := metrics.MetricDTOParams{MetricsName: name, MetricType: parameters[0]}
	p := []*metrics.MetricDTOParams{&metricParams}
	metric, err := metricsService.Get(storage.StorageInstance, &p)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)

		utils.WrireZeroBytes(res)
		return
	}
	res.Header().Set("Content-Type", " text/plain")
	if parameters[0] == constants.Gauge {
		_, err = res.Write([]byte(strconv.FormatFloat(*metric.Value, 'f', -1, 64)))
		if err != nil {
			logger.LogError(err)
			return
		}
		return
	}
	_, err = res.Write([]byte(strconv.FormatInt(int64(*metric.Delta), 10)))
	if err != nil {
		logger.LogError(err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// GetOneHandler handles HTTP POST requests to retrieve a single metric in JSON format.
// Accepts a metric object in the request body.
// Returns the current metric value as JSON.
// Responds with appropriate HTTP status codes for errors.
func GetOneHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("GetOneHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusMethodNotAllowed)

		utils.WrireZeroBytes(res)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusBadRequest)

		utils.WrireZeroBytes(res)
		return
	}
	requestMetric, err := parseMetric(&buf)
	logger.LogMetric(requestMetric)

	if err != nil {

		logger.LogError(err)
		res.WriteHeader(http.StatusBadRequest)

		utils.WrireZeroBytes(res)
		return
	}

	metricParams := metrics.MetricDTOParams{MetricsName: requestMetric.ID, MetricType: requestMetric.MType}
	p := []*metrics.MetricDTOParams{&metricParams}

	responseMetrics, err := metricsService.Get(storage.StorageInstance, &p)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		utils.WrireZeroBytes(res)
		return
	}

	body, err := json.Marshal(responseMetrics)

	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusInternalServerError)
		utils.WrireZeroBytes(res)
		return
	}
	logger.LogMetric(*responseMetrics)

	res.Header().Set("Content-Type", " application/json")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(body); err != nil {
		logger.LogError(err)
	}

}

// UpdateHandler handles HTTP POST requests to update a metric.
// Accepts a metric object in JSON format in the request body.
// Returns HTTP 200 on success.
// Responds with appropriate HTTP status codes for errors.
func UpdateHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("updateHandler")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogError(err)

		res.WriteHeader(http.StatusBadRequest)
		utils.WrireZeroBytes(res)
		return
	}

	metric, err := parseMetric(&buf)

	if err != nil {
		logger.LogError(err)
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		utils.WrireZeroBytes(res)
		return
	}

	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}

	res.WriteHeader(http.StatusOK)
	utils.WrireZeroBytes(res)
}

// UpdateHandlerByURLParams handles HTTP requests to update a metric via URL parameters.
// Expected URL format: /update/<type>/<name>/<value>.
// Returns HTTP 200 on success.
// Responds with appropriate HTTP status codes for errors.
func UpdateHandlerByURLParams(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("UpdateHandlerByURLParams \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost, http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}
	urlString := req.URL.Path //  /update/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/update/")
	parameters := strings.Split(params, "/")
	metric, err := metricRecord(parameters)

	if err != nil {
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		utils.WrireZeroBytes(res)
		return
	}
	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte(params)); err != nil {
		return
	}
	utils.WrireZeroBytes(res)
}

// UpdatesHandler handles HTTP POST requests for batch metric updates.
// Accepts an array of metric objects in JSON format.
// Returns HTTP 200 on success.
// Responds with appropriate HTTP status codes for errors.
func UpdatesHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("UpdatesHandler")

	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogError(err)

		res.WriteHeader(http.StatusBadRequest)
		utils.WrireZeroBytes(res)
		return
	}
	metricsSlice, err := parseMetrics(&buf)

	if err != nil {
		logger.LogError(err)
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		utils.WrireZeroBytes(res)
		return
	}
	err = metricsService.UpdateMany(storage.StorageInstance, metricsSlice)
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		utils.WrireZeroBytes(res)
		return
	}

	res.WriteHeader(http.StatusOK)
	utils.WrireZeroBytes(res)
}
func metricRecord(parameters []string) (metrics.Metrics, error) {
	if len(parameters) != 3 {
		return metrics.Metrics{}, ErrNoMetricName
	}
	if parameters[0] != constants.Gauge && parameters[0] != constants.Counter {
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
	if parameters[0] == constants.Gauge {
		value, err = strconv.ParseFloat(parameters[2], 64)
	} else {
		delta, err = strconv.ParseInt(parameters[2], 10, 64)

	}
	if err != nil {
		return metrics.Metrics{}, ErrWrongValue
	}

	return metrics.Metrics{
		MType: metricType,
		ID:    metricName,
		Value: &value,
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
		return metrics.Metrics{}, ErrNoMetricName
	}

	if metric.MType != constants.Gauge && metric.MType != constants.Counter {
		return metrics.Metrics{}, ErrNoMetricsType
	}
	if metric.ID == "" {
		return metrics.Metrics{}, ErrNoMetricName
	}
	return metric, nil
}
func parseMetrics(buf *bytes.Buffer) (*[]metrics.Metrics, error) {
	var metricsSlice []metrics.Metrics
	if err := json.Unmarshal(buf.Bytes(), &metricsSlice); err != nil {
		return &[]metrics.Metrics{}, ErrNoMetricName
	}
	for _, m := range metricsSlice {
		if m.MType != constants.Gauge && m.MType != constants.Counter {
			return &[]metrics.Metrics{}, ErrNoMetricsType
		}
		if m.ID == "" {
			return &[]metrics.Metrics{}, ErrNoMetricName
		}
	}

	return &metricsSlice, nil
}

// PingDB checks the database connection.
// Returns HTTP 200 if connection is successful,
// or HTTP 500 if there's a connection error.
func PingDB(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("PingDB")

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()
	err := storageService.Ping(ctx, storage.StorageInstance)
	if err != nil {
		logger.LogError(err)
		res.WriteHeader(http.StatusInternalServerError)
		utils.WrireZeroBytes(res)
		return
	}

	res.WriteHeader(http.StatusOK)
	utils.WrireZeroBytes(res)
}
