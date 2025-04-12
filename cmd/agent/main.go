package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
)


func main() {
	parameterts := config.GetParameters()

	httpClient := client.NewClient(parameterts.Addres)
	reportIntervalStart := time.Now()
	for {
		metrics, err := metricGenerator.Generator.Generate()
		if err != nil {
			logger.LogInfo(err)
			panic("Can not collect metrics")
		}
		if time.Since(reportIntervalStart) >= time.Duration(parameterts.ReportInterval)*time.Second {
			metricGenerator.Generator.UpdatePollCount()
			httpClient.SendMetrics(metrics)
			reportIntervalStart = time.Now()

		}
		time.Sleep(time.Duration(parameterts.PollInterval) * time.Second)
	}
}

