package metric

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/templates"
)
type Storage interface {
	SaveMetric(m *metrics.MetricDTO) error
	GetMetrics(metricsNames *[]string) (*[]metrics.MetricDTO , error)
}

func GetAll(s Storage) (string, error) {
	empySlice := []string{}
	metricsSlice, err := s.GetMetrics(&empySlice)
	if err != nil {
		return "", err
	}
	html := templates.GetAllMetricsHTMLPage(metricsSlice)
	return html, nil
}

func Get(s Storage, metricsNames *[]string) (*metrics.MetricDTO, error) {

	metricsSlice, err := s.GetMetrics(metricsNames)
	if err != nil {
		return nil, err
	}
	metric := (*metricsSlice)[0]
	return &metric, nil
}

func Update(s Storage, m *metrics.MetricDTO) error {
	err := s.SaveMetric(m)
	if err != nil {
		return err
	}
	return nil
}
