package generator

import (
	"math/rand"

	"runtime"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type MetricGenerator interface {
	Generate() ([]*metrics.Metrics, error)
	UpdatePollCount() int64
}

type generator struct {
	pollCount int64
}

var Generator = generator{
	pollCount: 0,
}

func (g *generator) Generate() ([]*metrics.Metrics, error) {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var metricSlice []*metrics.Metrics
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "Alloc", Value: utils.IntToPointerFloat(memStats.Alloc)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "BuckHashSys", Value: utils.IntToPointerFloat(memStats.BuckHashSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "Frees", Value: utils.IntToPointerFloat(memStats.Frees)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "GCCPUFraction", Value: utils.IntToPointerFloat(uint64(memStats.GCCPUFraction))})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "GCSys", Value: utils.IntToPointerFloat(memStats.GCSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapAlloc", Value: utils.IntToPointerFloat(memStats.HeapAlloc)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapIdle", Value: utils.IntToPointerFloat(memStats.HeapIdle)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapInuse", Value: utils.IntToPointerFloat(memStats.HeapInuse)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapObjects", Value: utils.IntToPointerFloat(memStats.HeapObjects)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapReleased", Value: utils.IntToPointerFloat(memStats.HeapReleased)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "HeapSys", Value: utils.IntToPointerFloat(memStats.HeapSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "LastGC", Value: utils.IntToPointerFloat(memStats.LastGC)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "Lookups", Value: utils.IntToPointerFloat(memStats.Lookups)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "MCacheInuse", Value: utils.IntToPointerFloat(memStats.MCacheInuse)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "MCacheSys", Value: utils.IntToPointerFloat(memStats.MCacheSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "MSpanInuse", Value: utils.IntToPointerFloat(memStats.MSpanInuse)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "MSpanSys", Value: utils.IntToPointerFloat(memStats.MSpanSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "Mallocs", Value: utils.IntToPointerFloat(memStats.Mallocs)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "NextGC", Value: utils.IntToPointerFloat(memStats.NextGC)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "NumForcedGC", Value: utils.IntToPointerFloat(uint64(memStats.NumForcedGC))})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "NumGC", Value: utils.IntToPointerFloat(uint64(memStats.NumGC))})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "OtherSys", Value: utils.IntToPointerFloat(memStats.OtherSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "PauseTotalNs", Value: utils.IntToPointerFloat(memStats.PauseTotalNs)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "StackInuse", Value: utils.IntToPointerFloat(memStats.StackInuse)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "StackSys", Value: utils.IntToPointerFloat(memStats.StackSys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "Sys", Value: utils.IntToPointerFloat(memStats.Sys)})
	metricSlice = append(metricSlice, &metrics.Metrics{MType: constants.Gauge, ID: "TotalAlloc", Value: utils.IntToPointerFloat(memStats.TotalAlloc)})

	metricSlice = append(metricSlice, &metrics.Metrics{
		MType: constants.Gauge,
		ID:    "RandomValue",
		Value: utils.IntToPointerFloat(rand.Uint64()),
	})

	metricSlice = append(metricSlice, &metrics.Metrics{
		MType: constants.Counter,
		ID:    "PollCount",
		Delta: (utils.IntToPointerInt(g.pollCount)),
	})
	logger.LogInfo(g.pollCount)
	return metricSlice, nil
}

func (g *generator) UpdatePollCount() int64 {
	g.pollCount++
	return g.pollCount
}
