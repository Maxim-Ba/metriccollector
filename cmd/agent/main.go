package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	metricgenerator "github.com/Maxim-Ba/metriccollector/internal/agent/metric_generator"
)



func main() {
	parseFlags()
	httpClient := client.NewClient(flagRunAddr)
	reportIntervalStart := time.Now()
	for {
		metrics, err := metricgenerator.Generator.Generate()
		if err != nil {
			panic("Can not collect metrics")
		}
		if time.Since(reportIntervalStart) >= time.Duration(flagReportInterval)* time.Second {
			metricgenerator.Generator.UpdatePollCount()
			httpClient.SendMetrics(metrics)
			reportIntervalStart = time.Now()

		}
		time.Sleep(time.Duration(flagPollInterval) * time.Second)
	}
}
