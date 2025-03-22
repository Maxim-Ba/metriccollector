package metricgenerator

import (
	"slices"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func Test_generator_updatePollCount(t *testing.T) {
	g := generator{pollCount: 5}

	result := g.updatePollCount()

	if g.pollCount != 6 {
		t.Errorf("Expected pollCount to be 6, got %d", g.pollCount)
	}

	if result != 6 {
		t.Errorf("Expected returned value to be 6, got %d", result)
	}
}

func Test_generator_Generate(t *testing.T) {
	g :=generator{pollCount: 0}
	result,err :=  g.Generate()
	if err!=nil {
		t.Error("Un expected error", err)
	}
	for _, v := range result {
		if !(slices.Contains(metrics.GaugeMetrics, v.MetricName))   {
			if v.MetricName != "PollCount"&& v.MetricName != "RandomValue" {
				t.Error("Un expected MetricName", v.MetricName)
			}
		}
	}
}
