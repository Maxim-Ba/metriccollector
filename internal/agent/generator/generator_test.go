package generator

import (
	"slices"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func Test_generator_UpdatePollCount(t *testing.T) {
	g := generator{pollCount: 5}

	result := g.UpdatePollCount()

	if g.pollCount != 6 {
		t.Errorf("Expected pollCount to be 6, got %d", g.pollCount)
	}

	if result != 6 {
		t.Errorf("Expected returned value to be 6, got %d", result)
	}
}

func Test_generator_Generate(t *testing.T) {
	g := generator{pollCount: 0}
	result, err := g.Generate(4)
	if err != nil {
		t.Error("Un expected error", err)
	}
	for _, v := range result {
		if !(slices.Contains(metrics.GaugeMetrics, v.ID)) {
			if v.ID != "PollCount" {
				t.Error("Un expected MetricName", v.ID)
			}
		}
	}
}
