package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
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
			err = utils.RetryWrapper(func() error {
				return httpClient.SendMetrics(metrics)
			}, 3, []error{client.ErrServerInternalError,client.ErrServerInternalError})
			// err := httpClient.SendMetrics(metrics)
			if err != nil {
				logger.LogError(err)
			}
			err = utils.RetryWrapper(func() error {
				return httpClient.SendMetricsWithBatch(metrics)
			}, 3, []error{client.ErrServerInternalError,client.ErrServerInternalError })
			// err = httpClient.SendMetricsWithBatch(metrics)
			if err != nil {
				logger.LogError(err)
			}

			reportIntervalStart = time.Now()

		}
		time.Sleep(time.Duration(parameters.PollInterval) * time.Second)
	}
}
