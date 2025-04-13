package storage

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type MemStorage struct {
	collectionGauge   map[string]float64
	collectionCounter map[string]int64
}

var StorageInstance = MemStorage{
	collectionGauge:   map[string]float64{},
	collectionCounter: map[string]int64{},
}

func New() (*MemStorage, error) {
	return &StorageInstance, nil
}
func (s MemStorage) SaveMetric(m *metrics.Metrics) error {

	if m.MType == "gauge" {
		StorageInstance.collectionGauge[m.ID] = *m.Value
	}
	if m.MType == "counter" {
		metricValue := int64(*m.Delta)
		if val, ok := StorageInstance.collectionCounter[m.ID]; ok {
			StorageInstance.collectionCounter[m.ID] = val + metricValue
		} else {
			StorageInstance.collectionCounter[m.ID] = metricValue
		}
	}
	return nil
}

// 

func (s MemStorage) GetMetrics(metricsParams *[]*metrics.MetricDTOParams) (*[]metrics.Metrics, error) {
	// func (s MemStorage) GetMetrics(metricsNames *[]string) (*[]metrics.Metrics , error){
	metricsNames := make([]string, len(*metricsParams))
	for i, m := range *metricsParams {
		metricsNames[i] = m.MetricsName
	}
	var metricsSlice []metrics.Metrics
	// Get all metrics
	if len(metricsNames) == 0 {
		for metric, value := range StorageInstance.collectionGauge {
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: metric, Value: utils.FloatToPointerFloat(value)})
		}
		for metric, value := range StorageInstance.collectionCounter {
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: metric, Delta: utils.FloatToPointerInt(value)})
		}
		return &metricsSlice, nil
	}
	for _, m := range *metricsParams {
		if m.MetricType == "gauge" {
			value, ok := StorageInstance.collectionGauge[m.MetricsName]
			if ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: m.MetricsName, Value: utils.FloatToPointerFloat(value)})
			} else {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: m.MetricsName, Value: utils.FloatToPointerFloat(0)})
			}
		}
		if m.MetricType == "counter" {
			value, ok := StorageInstance.collectionCounter[m.MetricsName]
			if ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: m.MetricsName, Delta: utils.FloatToPointerInt(value)})
			} else {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: m.MetricsName, Delta: utils.FloatToPointerInt(0)})
			}

		}
	}
	//Get choosen metrics
	// for metric, value := range StorageInstance.collectionGauge {
	// 	if slices.Contains(metricsNames,metric) {
	// 		metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: metric, Value: utils.FloatToPointerFloat(value)})
	// 	}
	// 	else {
	// 		metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: metric, Value: utils.FloatToPointerFloat(value)})
	// 	}
	// }
	// for metric, value := range StorageInstance.collectionCounter {
	// 	if slices.Contains(metricsNames,metric) {
	// 		metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: metric, Delta: utils.FloatToPointerInt(value)})
	// 	}
	// }
	if len(metricsSlice) == 0 {
		return nil, ErrUnknownMetricName
	}
	return &metricsSlice, nil
}
