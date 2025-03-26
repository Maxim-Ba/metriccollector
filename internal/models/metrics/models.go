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

var GaugeMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}
var CounterMetrics = []string{
	"RandomValue",
	"PollCount",
}
