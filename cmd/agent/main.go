package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	metricGenerator "github.com/Maxim-Ba/metriccollector/internal/agent/generator"
)

type Parameters struct {
	Addres         string
	ReportInterval int
	PollInterval   int
}

func main() {
	parameterts := getParameters()

	httpClient := client.NewClient(parameterts.Addres)
	reportIntervalStart := time.Now()
	for {
		metrics, err := metricGenerator.Generator.Generate()
		if err != nil {
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

func getParameters() Parameters {
	flags:=config.ParseFlags()
	envConfig := config.ParseEnv()
	address := envConfig.Address
	pollInterval := envConfig.PollInterval
	reportInterval := envConfig.ReportInterval
	if address == "" {
		address = flags.FlagRunAddr
	}
	if pollInterval == 0 {
		pollInterval = flags.FlagPollInterval
	}
	if reportInterval == 0 {
		reportInterval = flags.FlagReportInterval
	}
	return Parameters{
		Addres:         address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
