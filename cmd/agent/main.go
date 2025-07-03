package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

func main() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	parameters := config.New()
	logger.SetLogLevel(parameters.LogLevel)
	signature.New(parameters.Key)
	httpClient := client.NewClient(parameters.Addres)
	reportIntervalStart := time.Now()
	go func() {

		for {
			metrics, err := metricGenerator.Generator.Generate(parameters.RateLimit)
			if err != nil {
				logger.LogError(err)
				panic("Can not collect metrics")
			}
			if time.Since(reportIntervalStart) >= time.Duration(parameters.ReportInterval)*time.Second {
				metricGenerator.Generator.UpdatePollCount()
				err = utils.RetryWrapper(func() error {
					return httpClient.SendMetrics(metrics)
				}, []error{client.ErrServerInternalError, client.ErrRequestTimeout})
				if err != nil {
					logger.LogError(err)
				}
				err = utils.RetryWrapper(func() error {
					return httpClient.SendMetricsWithBatch(metrics)
				}, []error{client.ErrServerInternalError, client.ErrRequestTimeout})
				if err != nil {
					logger.LogError(err)
				}

				reportIntervalStart = time.Now()

			}
			time.Sleep(time.Duration(parameters.PollInterval) * time.Second)
		}
	}()

	<-exit // Ожидание сигнала завершения
	logger.LogInfo("Shutting down agent...")
	logger.Sync()

}
