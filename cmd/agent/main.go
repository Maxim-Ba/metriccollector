package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
	protoclient "github.com/Maxim-Ba/metriccollector/internal/agent/proto-client"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
	"github.com/Maxim-Ba/metriccollector/pkg/buildinfo"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	buildinfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	parameters := config.New()
	logger.SetLogLevel(parameters.LogLevel)
	signature.New(parameters.Key, parameters.CryptoKeyPath)
	httpClient := client.NewClient(parameters.Address)
	reportIntervalStart := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	conn, err := protoclient.Connect(parameters.GrpcServer)
	if err != nil && parameters.GrpcOn {
		panic(err)
	}
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				logger.LogInfo("Stopping metrics collection...")
				return
			default:
				metrics, err := metricGenerator.Generator.Generate(parameters.RateLimit)
				if err != nil {
					logger.LogError(err)
					panic("Can not collect metrics")
				}
				if time.Since(reportIntervalStart) >= time.Duration(parameters.ReportInterval)*time.Second {
					metricGenerator.Generator.UpdatePollCount()
					if parameters.GrpcOn {
						conn.SendMetrics(metrics)
					} else {

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
					}

					reportIntervalStart = time.Now()
				}
				time.Sleep(time.Duration(parameters.PollInterval) * time.Second)
			}
		}
	}()

	<-exit // Ожидание сигнала завершения
	cancel()
	wg.Wait()

	logger.LogInfo("Shutting down agent...")
	logger.Sync()
}
