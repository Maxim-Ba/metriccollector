package storage

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

type Storage interface {
	SaveMetric(m *metrics.MetricDTO) error
}


type MemStorage struct {
	collectionGauge map[string]float64
	collectionCounter map[string]int64
}
var storage = MemStorage{
	collectionGauge: map[string]float64{},
	collectionCounter: map[string]int64{},
}

func SaveMetric (m *metrics.MetricDTO) error {
	if m.MetricType == "gauge" {
		storage.collectionGauge[m.MetricName] = m.Value
	}
	if m.MetricType == "counter" {
		metricValue :=int64(m.Value) 
		if val,ok:= storage.collectionCounter[m.MetricName]; ok{
			storage.collectionCounter[m.MetricName] = val + metricValue
		} else {
			storage.collectionCounter[m.MetricName] = metricValue
		}
	}
return nil
}
