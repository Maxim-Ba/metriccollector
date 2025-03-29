package generator

import (
	"fmt"
	"math/rand"

	"runtime"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

type MetricGenerator interface {
	Generate() ([]*metrics.MetricDTO, error)
	UpdatePollCount() int64
}

type generator struct {
	pollCount int64
}



var Generator = generator{
	pollCount: 0,
}

func (g *generator) Generate() ([]*metrics.MetricDTO, error) {


	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var metricSlice []*metrics.MetricDTO
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "Alloc", Value: float64(memStats.Alloc)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "BuckHashSys", Value: float64(memStats.BuckHashSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "Frees", Value: float64(memStats.Frees)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "GCCPUFraction", Value: memStats.GCCPUFraction})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "GCSys", Value: float64(memStats.GCSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapAlloc", Value: float64(memStats.HeapAlloc)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapIdle", Value: float64(memStats.HeapIdle)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapInuse", Value: float64(memStats.HeapInuse)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapObjects", Value: float64(memStats.HeapObjects)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapReleased", Value: float64(memStats.HeapReleased)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "HeapSys", Value: float64(memStats.HeapSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "LastGC", Value: float64(memStats.LastGC)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "Lookups", Value: float64(memStats.Lookups)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "MCacheInuse", Value: float64(memStats.MCacheInuse)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "MCacheSys", Value: float64(memStats.MCacheSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "MSpanInuse", Value: float64(memStats.MSpanInuse)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "MSpanSys", Value: float64(memStats.MSpanSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "Mallocs", Value: float64(memStats.Mallocs)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "NextGC", Value: float64(memStats.NextGC)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "NumForcedGC", Value: float64(memStats.NumForcedGC)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "NumGC", Value: float64(memStats.NumGC)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "OtherSys", Value: float64(memStats.OtherSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "PauseTotalNs", Value: float64(memStats.PauseTotalNs)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "StackInuse", Value: float64(memStats.StackInuse)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "StackSys", Value: float64(memStats.StackSys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "Sys", Value: float64(memStats.Sys)})
	metricSlice = append(metricSlice, &metrics.MetricDTO{MetricType: "gauge", MetricName: "TotalAlloc", Value: float64(memStats.TotalAlloc)})



	metricSlice = append(metricSlice, &metrics.MetricDTO{
		MetricType: "counter",
		MetricName: "RandomValue",
		Value:      float64(rand.Int63()),
	})
	metricSlice = append(metricSlice, &metrics.MetricDTO{
		MetricType: "counter",
		MetricName: "PollCount",
		Value:      float64(g.pollCount),
	})
	fmt.Print(g.pollCount)
	return metricSlice, nil
}

func (g *generator) UpdatePollCount() int64 {
	g.pollCount = g.pollCount + 1
	return g.pollCount
}
