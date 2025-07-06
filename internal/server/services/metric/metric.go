package metric

import (
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/templates"
)

// Storage defines the interface for metric persistence operations.
// Implementations should provide methods for saving and retrieving metrics.
type Storage interface {
	SaveMetric(m *metrics.Metrics) error
	SaveMetrics(m *[]metrics.Metrics) error
	GetMetrics(params *[]*metrics.MetricDTOParams) (*[]metrics.Metrics, error)
}

// GetAll retrieves all metrics from storage and returns them as an HTML page.
func GetAll(s Storage) (string, error) {
	empySlice := []*metrics.MetricDTOParams{}
	metricsSlice, err := s.GetMetrics(&empySlice)
	if err != nil {
		return "", err
	}
	html := templates.GetAllMetricsHTMLPage(metricsSlice)
	return html, nil
}

// Get retrieves a specific metric from storage based on provided parameters.
func Get(s Storage, metricsNames *[]*metrics.MetricDTOParams) (*metrics.Metrics, error) {

	metricsSlice, err := s.GetMetrics(metricsNames)
	if err != nil {
		logger.LogInfo(err)
		return nil, err
	}
	metric := (*metricsSlice)[0]
	return &metric, nil
}

// Update persists a single metric to storage.
func Update(s Storage, m *metrics.Metrics) error {
	err := s.SaveMetric(m)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMany persists multiple metrics to storage in a batch operation.
func UpdateMany(s Storage, m *[]metrics.Metrics) error {
	err := s.SaveMetrics(m)
	if err != nil {
		return err
	}
	return nil
}
