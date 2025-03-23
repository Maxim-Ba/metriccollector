package main

import (
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/agent/client"
	metricgenerator "github.com/Maxim-Ba/metriccollector/internal/agent/metric_generator"
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
		metrics, err := metricgenerator.Generator.Generate()
		if err != nil {
			panic("Can not collect metrics")
		}
		if time.Since(reportIntervalStart) >= time.Duration(parameterts.ReportInterval)*time.Second {
			metricgenerator.Generator.UpdatePollCount()
			httpClient.SendMetrics(metrics)
			reportIntervalStart = time.Now()

		}
		time.Sleep(time.Duration(parameterts.PollInterval) * time.Second)
	}
}
func getParameters() Parameters {
	parseFlags()
	envConfig := parseEnv()
	address := envConfig.ADDRESS
	pollInterval := envConfig.POLL_INTERVAL
	reportInterval := envConfig.REPORT_INTERVAL
	if address == "" {
		address = flagRunAddr
	}
	if pollInterval == 0 {
		pollInterval = flagPollInterval
	}
	if reportInterval == 0 {
		reportInterval = flagReportInterval
	}
	return Parameters{
		Addres:         address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
