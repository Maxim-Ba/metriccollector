package metricgenerator

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

type MetricGenerator interface {
	Generate() ([]*metrics.MetricDTO, error)
	updatePollCount() int64
}

type generator struct {
	pollCount int64
}

var gaugeMetrics = []string{
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

var Generator = generator{
	pollCount: 0,
}

func (g generator) Generate() ([]*metrics.MetricDTO, error) {

	// var metricSlice []*metrics.MetricDTO
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "Alloc", float64(memStats.Alloc)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "BuckHashSys", float64(memStats.BuckHashSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "Frees", float64(memStats.Frees)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "GCCPUFraction", memStats.GCCPUFraction})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "GCSys", float64(memStats.GCSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapAlloc", float64(memStats.HeapAlloc)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapIdle", float64(memStats.HeapIdle)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapInuse", float64(memStats.HeapInuse)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapObjects", float64(memStats.HeapObjects)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapReleased", float64(memStats.HeapReleased)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "HeapSys", float64(memStats.HeapSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "LastGC", float64(memStats.LastGC)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "Lookups", float64(memStats.Lookups)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "MCacheInuse", float64(memStats.MCacheInuse)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "MCacheSys", float64(memStats.MCacheSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "MSpanInuse", float64(memStats.MSpanInuse)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "MSpanSys", float64(memStats.MSpanSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "Mallocs", float64(memStats.Mallocs)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "NextGC", float64(memStats.NextGC)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "NumForcedGC", float64(memStats.NumForcedGC)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "NumGC", float64(memStats.NumGC)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "OtherSys", float64(memStats.OtherSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "PauseTotalNs", float64(memStats.PauseTotalNs)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "StackInuse", float64(memStats.StackInuse)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "StackSys", float64(memStats.StackSys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "Sys", float64(memStats.Sys)})
	// metricSlice = append(metricSlice, &metrics.MetricDTO{"gauge", "TotalAlloc", float64(memStats.TotalAlloc)})
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var metricSlice []*metrics.MetricDTO
	memStatsValue := reflect.ValueOf(memStats)
	for _, metricName := range gaugeMetrics {
		field := memStatsValue.FieldByName(metricName)
		if field.IsValid() {
			var value float64
			switch field.Kind() {
			case reflect.Uint64:
				value = float64(field.Uint())
			case reflect.Float64:
				value = field.Float()
			}

			metricSlice = append(metricSlice, &metrics.MetricDTO{
				MetricType: "gauge",
				MetricName: metricName,
				Value:      value,
			})
		}

	}

	metricSlice = append(metricSlice, &metrics.MetricDTO{
		MetricType: "counter",
		MetricName: "RandomValue",
		Value:      float64(rand.Int63()),
	})
	metricSlice = append(metricSlice, &metrics.MetricDTO{
		MetricType: "counter",
		MetricName: "PollCount",
		Value:      float64(g.updatePollCount()),
	})

	return metricSlice, nil
}

func (g *generator) updatePollCount() int64 {
	g.pollCount = g.pollCount + 1
	return g.pollCount
}
