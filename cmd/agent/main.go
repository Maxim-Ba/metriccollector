package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	metricgenerator "github.com/Maxim-Ba/metriccollector/internal/agent/metric_generator"
)

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second

func main() {
	httpClient := client.NewClient("localhost", 8080)
	reportIntervalStart := time.Now()
	for {
		metrics, err := metricgenerator.Generator.Generate()
		if err != nil {
			panic("Can not collect metrics")
		}
		if time.Since(reportIntervalStart) >= reportInterval {
			httpClient.SendMetrics(metrics)
			reportIntervalStart = time.Now()

		}
		time.Sleep(pollInterval)
	}
}
