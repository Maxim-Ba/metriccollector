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
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

func GetAllHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("getAllHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		// _, err := res.Write([]byte(""))
		// if err != nil {
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}

	html, err := metricsService.GetAll(storage.StorageInstance)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		// _, err := res.Write([]byte(""))
		// if err != nil {
		// 	return
		// }
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
func GetOneHandlerByParams(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("GetOneHandlerByParams")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		// _, err := res.Write([]byte(""))
		// if err != nil {
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}
	urlString := req.URL.Path //  /value/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/value/")
	parameters := strings.Split(params, "/")
	name := parameters[1]
	metricParams := metrics.MetricDTOParams{MetricsName: name, MetricType: parameters[0]}
	p := []*metrics.MetricDTOParams{&metricParams}
	logger.LogInfo(parameters[1]) //
	metric, err := metricsService.Get(storage.StorageInstance, &p)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		// _, err := res.Write([]byte(""))
		// if err != nil {
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}
	res.Header().Set("Content-Type", " text/plain")
	if parameters[0] == "gauge" {
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
func GetOneHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("GetOneHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		// if _, err := res.Write([]byte("")); err != nil {
		// 	logger.LogError(err)
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusBadRequest)
		// if _, err := res.Write([]byte("")); err != nil {
		// 	logger.LogError(err)
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}
	requestMetric, err := parseMetric(&buf)
	logger.LogMetric(requestMetric)

	if err != nil {

		logger.LogInfo(err)
		res.WriteHeader(http.StatusBadRequest)
		// if _, err := res.Write([]byte("")); err != nil {
		// 	logger.LogError(err)
		// 	return
		// }
		utils.WrireZeroBytes(res)
		return
	}

	// metricsNames :=[]string{requestMetric.ID}

	metricParams := metrics.MetricDTOParams{MetricsName: requestMetric.ID, MetricType: requestMetric.MType}
	p := []*metrics.MetricDTOParams{&metricParams}

	responseMetrics, err := metricsService.Get(storage.StorageInstance, &p)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}

	body, err := json.Marshal(responseMetrics)

	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusInternalServerError)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}
	logger.LogMetric(*responseMetrics)

	logger.LogInfo("return GetOneHandler")

	res.Header().Set("Content-Type", " application/json")
	res.WriteHeader(http.StatusOK)
	if _,err:= res.Write(body); err != nil {
		logger.LogError(err)
	}
	
}

func UpdateHandler(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("updateHandler \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		logger.LogInfo(err)

		res.WriteHeader(http.StatusBadRequest)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}

	metric, err := parseMetric(&buf)

	if err != nil {
		logger.LogInfo(err)
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}
	// logger.LogMetric(metric)

	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}

	res.WriteHeader(http.StatusOK)
	// res.Write([]byte(""))
	utils.WrireZeroBytes(res)
}
func UpdateHandlerByURLParams(res http.ResponseWriter, req *http.Request) {
	logger.LogInfo("UpdateHandlerByURLParams \n")
	err := checkForAllowedMethod(req, []string{http.MethodPost, http.MethodGet})
	if err != nil {
		// logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}
	urlString := req.URL.Path //  /update/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/update/")
	parameters := strings.Split(params, "/")
	metric, err := metricRecord(parameters)

	if err != nil {
		// logger.LogInfo(err)
		if err == ErrNoMetricName {
			res.WriteHeader(http.StatusNotFound)
		}
		if err == ErrNoMetricsType || err == ErrWrongValue {
			res.WriteHeader(http.StatusBadRequest)
		}
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}
	err = metricsService.Update(storage.StorageInstance, &metric)
	if err != nil {
		// logger.LogInfo(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		// res.Write([]byte(""))
		utils.WrireZeroBytes(res)
		return
	}
	res.WriteHeader(http.StatusOK)
	if _,err:= res.Write([]byte(params));err != nil {
		return 
	} 
	// res.Write([]byte(""))
	utils.WrireZeroBytes(res)
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
		// logger.LogInfo(err)
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
