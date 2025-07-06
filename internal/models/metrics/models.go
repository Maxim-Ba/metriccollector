package metrics

// MetricDTO represents a Data Transfer Object for metric values.
// It contains the metric type, name, and current value.
// Used for internal processing and storage operations.
type MetricDTO struct {
	MetricType string
	MetricName string
	Value      float64
}

// Metrics represents a metric entity in JSON format.
// Used for API communication and serialization/deserialization.
// Fields:
//   - ID: Metric name (e.g., "HeapAlloc")
//   - MType: Metric type, either "gauge" or "counter"
//   - Delta: Pointer to integer value for counter metrics (optional)
//   - Value: Pointer to float value for gauge metrics (optional)
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// MetricDTOParams contains parameters for metric lookup operations.
// Used when querying metrics by name and type.
type MetricDTOParams struct {
	MetricsName string
	MetricType  string
}

// GaugeMetrics contains all supported gauge metric names.
// These represent runtime metrics that can increase or decrease.
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
	"TotalAlloc", "TotalMemory", "FreeMemory", "CPUutilization1", "RandomValue",
}

// CounterMetrics contains all supported counter metric names.
// These represent monotonically increasing counters.
var CounterMetrics = []string{
	"RandomValue",
	"PollCount",
}
