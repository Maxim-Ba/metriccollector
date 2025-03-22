package storage

import (
	"slices"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

type Storage interface {
	SaveMetric(m *metrics.MetricDTO) error
	GetMetrics(metricsNames []string) *[]metrics.MetricDTO
}

type MemStorage struct {
	collectionGauge   map[string]float64
	collectionCounter map[string]int64
}

var storage = MemStorage{
	collectionGauge:   map[string]float64{},
	collectionCounter: map[string]int64{},
}

func SaveMetric(m *metrics.MetricDTO) error {
	if m.MetricType == "gauge" {
		storage.collectionGauge[m.MetricName] = m.Value
	}
	if m.MetricType == "counter" {
		metricValue := int64(m.Value)
		if val, ok := storage.collectionCounter[m.MetricName]; ok {
			storage.collectionCounter[m.MetricName] = val + metricValue
		} else {
			storage.collectionCounter[m.MetricName] = metricValue
		}
	}
	return nil
}
func GetMetrics(metricsNames *[]string) (*[]metrics.MetricDTO , error){
	err:= checkCorrectMetricsName(metricsNames)
	if err != nil {
		return nil, ErrUnknownMetricName
	}
	var metricsSlice []metrics.MetricDTO
	// Get all metrics
	if len(*metricsNames) == 0 {
		for metric, value := range storage.collectionGauge {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "gauge", MetricName: metric, Value: value})
		}
		for metric, value := range storage.collectionCounter {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "counter", MetricName: metric, Value: float64( value)})
		}
		return &metricsSlice, nil
	}
	//Get choosen metrics
	for metric, value := range storage.collectionGauge {
		if slices.Contains(*metricsNames,metric) {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "gauge", MetricName: metric, Value: value})
		}
	}
	for metric, value := range storage.collectionCounter {
		if slices.Contains(*metricsNames,metric) {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "counter", MetricName: metric, Value: float64( value)})
		}
	}
	
	return &metricsSlice, nil
}

func checkCorrectMetricsName(metricsNames *[]string) error {
	for _, name := range *metricsNames {
		if !(slices.Contains(metrics.GaugeMetrics, name)) && !(slices.Contains(metrics.CounterMetrics, name)) {
			return ErrUnknownMetricName
		}
	}
	return nil
}
