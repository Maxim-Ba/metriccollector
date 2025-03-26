package storage

import "github.com/Maxim-Ba/metriccollector/internal/models/metrics"

type Storage interface {
	SaveMetric(m *metrics.MetricDTO) error
	GetMetrics(metricsNames *[]string) (*[]metrics.MetricDTO , error)
}
