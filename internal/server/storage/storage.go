package storage

import (
	"slices"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)


type MemStorage struct {
	collectionGauge   map[string]float64
	collectionCounter map[string]int64
}

var StorageInstance = MemStorage{
	collectionGauge:   map[string]float64{},
	collectionCounter: map[string]int64{},
}

func New()(*MemStorage ,error)  {
	return &StorageInstance, nil
}
func (s MemStorage) SaveMetric(m *metrics.MetricDTO) error {
	if m.MetricType == "gauge" {
		StorageInstance.collectionGauge[m.MetricName] = m.Value
	}
	if m.MetricType == "counter" {
		metricValue := int64(m.Value)
		if val, ok := StorageInstance.collectionCounter[m.MetricName]; ok {
			StorageInstance.collectionCounter[m.MetricName] = val + metricValue
		} else {
			StorageInstance.collectionCounter[m.MetricName] = metricValue
		}
	}
	return nil
}
func (s MemStorage) GetMetrics(metricsNames *[]string) (*[]metrics.MetricDTO , error){

	var metricsSlice []metrics.MetricDTO
	// Get all metrics
	if len(*metricsNames) == 0 {
		for metric, value := range StorageInstance.collectionGauge {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "gauge", MetricName: metric, Value: value})
		}
		for metric, value := range StorageInstance.collectionCounter {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "counter", MetricName: metric, Value: float64( value)})
		}
		return &metricsSlice, nil
	}
	//Get choosen metrics
	for metric, value := range StorageInstance.collectionGauge {
		if slices.Contains(*metricsNames,metric) {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "gauge", MetricName: metric, Value: value})
		}
	}
	for metric, value := range StorageInstance.collectionCounter {
		if slices.Contains(*metricsNames,metric) {
			metricsSlice = append(metricsSlice, metrics.MetricDTO{MetricType: "counter", MetricName: metric, Value: float64( value)})
		}
	}
	if len(metricsSlice) ==0 {
		return  nil, ErrUnknownMetricName
	}
	return &metricsSlice, nil
}
