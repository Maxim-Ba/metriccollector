package generator

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

func (g *generator) Generate(maxWorkers int) ([]*metrics.Metrics, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var metricSlice []*metrics.Metrics

	callbacks := []func() *metrics.Metrics{
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "Alloc",
				Value: utils.IntToPointerFloat(memStats.Alloc),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "BuckHashSys",
				Value: utils.IntToPointerFloat(memStats.BuckHashSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "Frees",
				Value: utils.IntToPointerFloat(memStats.Frees),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "GCCPUFraction",
				Value: utils.IntToPointerFloat(uint64(memStats.GCCPUFraction)),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "GCSys",
				Value: utils.IntToPointerFloat(memStats.GCSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapAlloc",
				Value: utils.IntToPointerFloat(memStats.HeapAlloc),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapIdle",
				Value: utils.IntToPointerFloat(memStats.HeapIdle),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapInuse",
				Value: utils.IntToPointerFloat(memStats.HeapInuse),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapObjects",
				Value: utils.IntToPointerFloat(memStats.HeapObjects),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapReleased",
				Value: utils.IntToPointerFloat(memStats.HeapReleased),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "HeapSys",
				Value: utils.IntToPointerFloat(memStats.HeapSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "LastGC",
				Value: utils.IntToPointerFloat(memStats.LastGC),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "Lookups",
				Value: utils.IntToPointerFloat(memStats.Lookups),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "MCacheInuse",
				Value: utils.IntToPointerFloat(memStats.MCacheInuse),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "MCacheSys",
				Value: utils.IntToPointerFloat(memStats.MCacheSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "MSpanInuse",
				Value: utils.IntToPointerFloat(memStats.MSpanInuse),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "MSpanSys",
				Value: utils.IntToPointerFloat(memStats.MSpanSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "Mallocs",
				Value: utils.IntToPointerFloat(memStats.Mallocs),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "NextGC",
				Value: utils.IntToPointerFloat(memStats.NextGC),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "NumForcedGC",
				Value: utils.IntToPointerFloat(uint64(memStats.NumForcedGC)),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "NumGC",
				Value: utils.IntToPointerFloat(uint64(memStats.NumGC)),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "OtherSys",
				Value: utils.IntToPointerFloat(memStats.OtherSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "PauseTotalNs",
				Value: utils.IntToPointerFloat(memStats.PauseTotalNs),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "StackInuse",
				Value: utils.IntToPointerFloat(memStats.StackInuse),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "StackSys",
				Value: utils.IntToPointerFloat(memStats.StackSys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "Sys",
				Value: utils.IntToPointerFloat(memStats.Sys),
			}
		},
		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "TotalAlloc",
				Value: utils.IntToPointerFloat(memStats.TotalAlloc),
			}
		},

		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Gauge,
				ID:    "RandomValue",
				Value: utils.IntToPointerFloat(rand.Uint64()),
			}
		},

		func() *metrics.Metrics {
			return &metrics.Metrics{
				MType: constants.Counter,
				ID:    "PollCount",
				Delta: (utils.IntToPointerInt(g.pollCount)),
			}
		},
	}

	doneCh := make(chan struct{})
	wpResultCh := workerPool(callbacks, maxWorkers)

	gopsutilResultCh := getGopsutilMetrics()

	resultCh := fanIn(doneCh, wpResultCh, gopsutilResultCh)
	for m := range resultCh {
		metricSlice = append(metricSlice, m)
	}
	return metricSlice, nil
}

func (g *generator) UpdatePollCount() int64 {
	g.pollCount++
	return g.pollCount
}

func getGopsutilMetrics() <-chan *metrics.Metrics {
	result := make(chan *metrics.Metrics, 3)
	go func() {
		defer close(result) // Закрываем канал после отправки всех данных
		v, err := mem.VirtualMemory()
		if err != nil {
			logger.LogError(err)
		}

		cpuUtilization, err := cpu.Percent(1*time.Second, true)
		if err != nil {
			logger.LogError(err)
			return
		}

		// Средняя загрузка всех ядер
		if len(cpuUtilization) == 0 {
			logger.LogError("cpuUtilization == 0")
		}
		totalUtil := 0.0
		for _, util := range cpuUtilization {
			totalUtil += util
		}
		averageUtil := totalUtil / float64(len(cpuUtilization))

		result <- &metrics.Metrics{
			MType: constants.Gauge,
			ID:    "TotalMemory",
			Value: utils.IntToPointerFloat(v.Total),
		}
		result <- &metrics.Metrics{
			MType: constants.Gauge,
			ID:    "FreeMemory",
			Value: utils.IntToPointerFloat(v.Free),
		}
		result <- &metrics.Metrics{
			MType: constants.Gauge,
			ID:    "CPUutilization1",
			Value: &averageUtil,
		}
	}()

	return result
}
