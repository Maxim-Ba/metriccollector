package storage

import (
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

func TestGetMetrics(t *testing.T) {
	type test struct {
		name     string
		params   *[]*metrics.MetricDTOParams
		want     *[]metrics.Metrics
		wantErr  bool
		fixtures MemStorage
	}
	tests := []test{
		{
			name: "happy pass gauge",
			params: &[]*metrics.MetricDTOParams{{MetricsName: "a",
				MetricType: "gauge"},
			},
			fixtures: MemStorage{
				collectionCounter: map[string]int64{
					string("a"): 100,
				},
				collectionGauge: map[string]float64{
					string("a"): 100,
				},
			},
			want: &[]metrics.Metrics{
				{ID: "a", MType: "gauge", Value: utils.IntToPointerFloat(100)},
			},
			wantErr: false,
		},
		{
			name: "happy pass counter",
			params: &[]*metrics.MetricDTOParams{{MetricsName: "a",
				MetricType: "counter"},
			},
			fixtures: MemStorage{
				collectionCounter: map[string]int64{
					string("a"): 100,
				},
				collectionGauge: map[string]float64{
					string("a"): 100,
				},
			},
			want: &[]metrics.Metrics{
				{ID: "a", MType: "counter", Delta: utils.IntToPointerInt(100)},
			},
			wantErr: false,
		},
		{
			name: "happy pass both types",
			params: &[]*metrics.MetricDTOParams{{MetricsName: "a",
				MetricType: "counter"}, {MetricsName: "random",
				MetricType: "counter"},{MetricsName: "a",
				MetricType: "gauge"},{MetricsName: "b",
				MetricType: "gauge"},
			},
			fixtures: MemStorage{
				collectionCounter: map[string]int64{
					string("a"): 100,
					string("random"): 100,
				},
				collectionGauge: map[string]float64{
					string("a"): 100,
					string("b"): 100,
				},
			},
			want: &[]metrics.Metrics{
				{ID: "a", MType: "counter", Delta: utils.IntToPointerInt(100)},
				{ID: "random", MType: "counter", Delta: utils.IntToPointerInt(100)},
				{ID: "a", MType: "gauge", Value: utils.IntToPointerFloat(100)},
				{ID: "b", MType: "gauge", Value: utils.IntToPointerFloat(100)},
			},
			wantErr: false,
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.collectionCounter = tt.fixtures.collectionCounter
			s.collectionGauge = tt.fixtures.collectionGauge
			got, err := s.GetMetrics(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("log %s, got %v", tt.name, got)
			gotSlice := *got
			for i, v := range *tt.want {
				if v.ID != gotSlice[i].ID {
					t.Errorf("GetMetrics() = %v, want %v", got, tt.want)
				}
				if v.MType != gotSlice[i].MType {
					t.Errorf("GetMetrics() = %v, want %v", got, tt.want)
				}
				if v.MType == "gauge" {
					if v.Value == nil || gotSlice[i].Value == nil || *v.Value != *(gotSlice[i].Value) {
						t.Errorf("GetMetrics() Value = %v, want %v", gotSlice[i].Value, v.Value)
					}
				}
				if v.MType == "counter" {
					if v.Delta == nil || gotSlice[i].Delta == nil || *v.Delta != *(gotSlice[i].Delta) {
						t.Errorf("GetMetrics() Delta = %v, want %v", gotSlice[i].Delta, v.Delta)
					}
				}
				
			}
			
		})
	}
}
