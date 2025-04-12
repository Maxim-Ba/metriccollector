package metric

import (
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/templates"
)
type Storage interface {
	SaveMetric(m *metrics.Metrics) error
	GetMetrics(metricsNames *[]string) (*[]metrics.Metrics , error)
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

func Get(s Storage, metricsNames *[]string) (*metrics.Metrics, error) {

	metricsSlice, err := s.GetMetrics(metricsNames)
	if err != nil {
		logger.LogInfo(err)
		return nil, err
	}
	metric := (*metricsSlice)[0]
	return &metric, nil
}

func Update(s Storage, m *metrics.Metrics) error {
	err := s.SaveMetric(m)
	if err != nil {
		return err
	}
	return nil
}
