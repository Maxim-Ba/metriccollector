package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"github.com/Maxim-Ba/metriccollector/internal/templates"
	"github.com/go-chi/chi/v5"
)

func InitHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", getAllHandler)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", getOneHandler)
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{value}", updateHandler)
	})

	return r
}
func getAllHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getAllHandler")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	empySlice := []string{}
	metricsSlice, er := storage.GetMetrics(&empySlice)
	if er != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(""))
		return
	}
	html := templates.GetAllMetricsHTMLPage(metricsSlice)
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(html))
}
func getOneHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("getOneHandler")
	err := checkForAllowedMethod(req, []string{http.MethodGet})
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(""))
		return
	}
	urlString := req.URL.Path //  /value/asdasd/asdasd/sdfsdfsdf/234
	params := strings.TrimPrefix(urlString, "/value/")
	parameters := strings.Split(params, "/")
	name := []string{parameters[1]}

	metricsSlice, err := storage.GetMetrics(&name)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(""))
		return
	}
	metric := (*metricsSlice)[0]
	res.Header().Set("Content-Type", " text/plain")
	res.Write([]byte(strconv.FormatFloat(metric.Value, 'f', 3, 64)))
}

func updateHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Print("updateHandler")
	err := checkForAllowedMethod(req, []string{http.MethodPost})
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
		if err == ErrNoMetricsType || err == ErrWrongValue {
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

func checkForAllowedMethod(req *http.Request, allowedMethod []string) error {
	if !(slices.Contains(allowedMethod, req.Method)) {
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
	if parameters[1] == "" {
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
