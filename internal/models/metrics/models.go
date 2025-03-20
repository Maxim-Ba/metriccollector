package metrics

type MetricDTO struct {
	MetricType string
	MetricName string
	Value      float64
}
type MetricDTOGauge struct {
	MetricDTO
	value float64
}
type MetricDTOCounter struct {
	MetricDTO
	value int64
}
