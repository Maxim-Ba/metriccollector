package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

func main() {
	defer logger.Sync()

	parameters := config.GetParameters()
	logger.SetLogLevel(parameters.LogLevel)
	httpClient := client.NewClient(parameters.Addres)
	reportIntervalStart := time.Now()
	for {
		metrics, err := metricGenerator.Generator.Generate()
		if err != nil {
			logger.LogError(err)
			panic("Can not collect metrics")
		}
		if time.Since(reportIntervalStart) >= time.Duration(parameters.ReportInterval)*time.Second {
			metricGenerator.Generator.UpdatePollCount()
			err := httpClient.SendMetrics(metrics)
			if err != nil {
				logger.LogError(err)
			}
			reportIntervalStart = time.Now()

		}
		time.Sleep(time.Duration(parameters.PollInterval) * time.Second)
	}
}
