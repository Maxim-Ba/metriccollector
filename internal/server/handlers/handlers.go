package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)

func InitHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateHandler)
	return mux
}

func updateHandler(res http.ResponseWriter, req *http.Request) {
	
	err := checkForAllowedMethod(req)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
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
		if err == ErrNoMetricsType||err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		res.Write([]byte(""))
		return
	}
	err = storage.SaveMetric(&metric)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(params))
}

func checkForAllowedMethod(req *http.Request) error {
	if req.Method != http.MethodPost {
		return fmt.Errorf("not allowed method")
	}
	return nil
}
func metricRecord(parameters []string) (metrics.MetricDTO, error) {
	if len(parameters) != 3 {
		return metrics.MetricDTO{}, ErrNoMetricName
	}
	if parameters[0] != "gauge" && parameters[0] != "counter" {
		return metrics.MetricDTO{}, ErrNoMetricsType
	}
	if parameters[1] == ""  {
		return metrics.MetricDTO{}, ErrNoMetricName
	}
	metricType := parameters[0]
	metricName := parameters[1]
	value, err := strconv.ParseFloat(parameters[2], 64)
	if err != nil {
		return metrics.MetricDTO{}, ErrWrongValue
	}
	return metrics.MetricDTO{
		MetricType: metricType,
		MetricName: metricName,
		Value:      value,
	}, nil
}
