package storage

import (
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
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
				MetricType: constants.Gauge},
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
				{ID: "a", MType: constants.Gauge, Value: utils.IntToPointerFloat(100)},
			},
			wantErr: false,
		},
		{
			name: "happy pass counter",
			params: &[]*metrics.MetricDTOParams{{MetricsName: "a",
				MetricType: constants.Counter},
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
				{ID: "a", MType: constants.Counter, Delta: utils.IntToPointerInt(100)},
			},
			wantErr: false,
		},
		{
			name: "happy pass both types",
			params: &[]*metrics.MetricDTOParams{{MetricsName: "a",
				MetricType: constants.Counter}, {MetricsName: "random",
				MetricType: constants.Counter}, {MetricsName: "a",
				MetricType: constants.Gauge}, {MetricsName: "b",
				MetricType: constants.Gauge},
			},
			fixtures: MemStorage{
				collectionCounter: map[string]int64{
					string("a"):      100,
					string("random"): 100,
				},
				collectionGauge: map[string]float64{
					string("a"): 100,
					string("b"): 100,
				},
			},
			want: &[]metrics.Metrics{
				{ID: "a", MType: constants.Counter, Delta: utils.IntToPointerInt(100)},
				{ID: "random", MType: constants.Counter, Delta: utils.IntToPointerInt(100)},
				{ID: "a", MType: constants.Gauge, Value: utils.IntToPointerFloat(100)},
				{ID: "b", MType: constants.Gauge, Value: utils.IntToPointerFloat(100)},
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
				if v.MType == constants.Gauge {
					if v.Value == nil || gotSlice[i].Value == nil || *v.Value != *(gotSlice[i].Value) {
						t.Errorf("GetMetrics() Value = %v, want %v", gotSlice[i].Value, v.Value)
					}
				}
				if v.MType == constants.Counter {
					if v.Delta == nil || gotSlice[i].Delta == nil || *v.Delta != *(gotSlice[i].Delta) {
						t.Errorf("GetMetrics() Delta = %v, want %v", gotSlice[i].Delta, v.Delta)
					}
				}

			}

		})
	}
}

func TestSaveMetric(t *testing.T) {
	type args struct {
		metric *metrics.Metrics
	}
	tests := []struct {
		name           string
		args           args
		wantGaugeVal   float64
		wantCounterVal int64
		wantErr        bool
		preSetup       func(s *MemStorage)
	}{
		{
			name: "save new gauge metric",
			args: args{
				metric: &metrics.Metrics{
					ID:    "test_gauge",
					MType: constants.Gauge,
					Value: utils.FloatToPointerFloat(123.45),
				},
			},
			wantGaugeVal: 123.45,
			wantErr:      false,
		},
		{
			name: "update existing gauge metric",
			args: args{
				metric: &metrics.Metrics{
					ID:    "existing_gauge",
					MType: constants.Gauge,
					Value: utils.FloatToPointerFloat(200.0),
				},
			},
			preSetup: func(s *MemStorage) {
				s.collectionGauge["existing_gauge"] = 100.0
			},
			wantGaugeVal: 200.0,
			wantErr:      false,
		},
		{
			name: "save new counter metric",
			args: args{
				metric: &metrics.Metrics{
					ID:    "test_counter",
					MType: constants.Counter,
					Delta: utils.IntToPointerInt(10),
				},
			},
			wantCounterVal: 10,
			wantErr:        false,
		},
		{
			name: "increment existing counter metric",
			args: args{
				metric: &metrics.Metrics{
					ID:    "existing_counter",
					MType: constants.Counter,
					Delta: utils.IntToPointerInt(5),
				},
			},
			preSetup: func(s *MemStorage) {
				s.collectionCounter["existing_counter"] = 10
			},
			wantCounterVal: 15,
			wantErr:        false,
		},
		{
			name: "invalid metric type",
			args: args{
				metric: &metrics.Metrics{
					ID:    "invalid",
					MType: "invalid_type",
				},
			},
			wantErr: true,
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ClearAll()

			if tt.preSetup != nil {
				tt.preSetup(s)
			}

			err := s.SaveMetric(tt.args.metric)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.args.metric.MType == constants.Gauge && !tt.wantErr {
				if val, ok := s.collectionGauge[tt.args.metric.ID]; !ok || val != tt.wantGaugeVal {
					t.Errorf("SaveMetric() gauge value = %v, want %v", val, tt.wantGaugeVal)
				}
			}

			if tt.args.metric.MType == constants.Counter && !tt.wantErr {
				if val, ok := s.collectionCounter[tt.args.metric.ID]; !ok || val != tt.wantCounterVal {
					t.Errorf("SaveMetric() counter value = %v, want %v", val, tt.wantCounterVal)
				}
			}
		})
	}
}

func TestSaveMetrics(t *testing.T) {
	type args struct {
		metrics *[]metrics.Metrics
	}
	tests := []struct {
		name            string
		args            args
		wantGaugeVals   map[string]float64
		wantCounterVals map[string]int64
		wantErr         bool
		preSetup        func(s *MemStorage)
	}{
		{
			name: "save multiple new metrics",
			args: args{
				metrics: &[]metrics.Metrics{
					{
						ID:    "gauge1",
						MType: constants.Gauge,
						Value: utils.FloatToPointerFloat(10.5),
					},
					{
						ID:    "counter1",
						MType: constants.Counter,
						Delta: utils.IntToPointerInt(5),
					},
					{
						ID:    "gauge2",
						MType: constants.Gauge,
						Value: utils.FloatToPointerFloat(20.0),
					},
				},
			},
			wantGaugeVals: map[string]float64{
				"gauge1": 10.5,
				"gauge2": 20.0,
			},
			wantCounterVals: map[string]int64{
				"counter1": 5,
			},
			wantErr: false,
		},
		{
			name: "update existing metrics",
			args: args{
				metrics: &[]metrics.Metrics{
					{
						ID:    "existing_gauge",
						MType: constants.Gauge,
						Value: utils.FloatToPointerFloat(30.0),
					},
					{
						ID:    "existing_counter",
						MType: constants.Counter,
						Delta: utils.IntToPointerInt(10),
					},
				},
			},
			preSetup: func(s *MemStorage) {
				s.collectionGauge["existing_gauge"] = 15.0
				s.collectionCounter["existing_counter"] = 5
			},
			wantGaugeVals: map[string]float64{
				"existing_gauge": 30.0,
			},
			wantCounterVals: map[string]int64{
				"existing_counter": 15,
			},
			wantErr: false,
		},
		{
			name: "mixed new and existing metrics",
			args: args{
				metrics: &[]metrics.Metrics{
					{
						ID:    "new_gauge",
						MType: constants.Gauge,
						Value: utils.FloatToPointerFloat(100.0),
					},
					{
						ID:    "existing_counter",
						MType: constants.Counter,
						Delta: utils.IntToPointerInt(3),
					},
				},
			},
			preSetup: func(s *MemStorage) {
				s.collectionCounter["existing_counter"] = 7
			},
			wantGaugeVals: map[string]float64{
				"new_gauge": 100.0,
			},
			wantCounterVals: map[string]int64{
				"existing_counter": 10,
			},
			wantErr: false,
		},
		{
			name: "invalid metric type",
			args: args{
				metrics: &[]metrics.Metrics{
					{
						ID:    "invalid",
						MType: "invalid_type",
					},
				},
			},
			wantErr: true,
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ClearAll()

			if tt.preSetup != nil {
				tt.preSetup(s)
			}

			err := s.SaveMetrics(tt.args.metrics)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for name, wantVal := range tt.wantGaugeVals {
					if val, ok := s.collectionGauge[name]; !ok || val != wantVal {
						t.Errorf("SaveMetrics() gauge %s value = %v, want %v", name, val, wantVal)
					}
				}

				for name, wantVal := range tt.wantCounterVals {
					if val, ok := s.collectionCounter[name]; !ok || val != wantVal {
						t.Errorf("SaveMetrics() counter %s value = %v, want %v", name, val, wantVal)
					}
				}
			}
		})
	}
}

func TestClearGaugeMetric(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		preSetup func(s *MemStorage)
		want     float64
		wantOk   bool
	}{
		{
			name: "clear existing gauge metric",
			args: args{
				name: "existing_gauge",
			},
			preSetup: func(s *MemStorage) {
				s.collectionGauge["existing_gauge"] = 100.0
			},
			want:   0,
			wantOk: false,
		},

		{
			name: "clear gauge from empty storage",
			args: args{
				name: "any_gauge",
			},
			preSetup: func(s *MemStorage) {
			},
			want:   0,
			wantOk: false,
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ClearAll()

			if tt.preSetup != nil {
				tt.preSetup(s)
			}

			s.ClearGaugeMetric(tt.args.name)

			val, ok := s.collectionGauge[tt.args.name]
			if ok != tt.wantOk {
				t.Errorf("ClearGaugeMetric() metric %s presence = %v, want %v", tt.args.name, ok, tt.wantOk)
			}
			if val != tt.want {
				t.Errorf("ClearGaugeMetric() metric %s value = %v, want %v", tt.args.name, val, tt.want)
			}

			if tt.name == "clear non-existing gauge metric" {
				if otherVal, otherOk := s.collectionGauge["other_gauge"]; !otherOk || otherVal != tt.want {
					t.Errorf("ClearGaugeMetric() affected other metric 'other_gauge' = %v, want %v", otherVal, tt.want)
				}
			}
		})
	}
}

func TestClearCounterMetric(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		preSetup func(s *MemStorage)
		want     int64
		wantOk   bool
	}{
		{
			name: "clear existing counter metric",
			args: args{
				name: "existing_counter",
			},
			preSetup: func(s *MemStorage) {
				s.collectionCounter["existing_counter"] = 100
			},
			want:   0,
			wantOk: false,
		},

		{
			name: "clear counter from empty storage",
			args: args{
				name: "any_counter",
			},
			preSetup: func(s *MemStorage) {
			},
			want:   0,
			wantOk: false,
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ClearAll()

			if tt.preSetup != nil {
				tt.preSetup(s)
			}

			s.ClearCounterMetric(tt.args.name)

			val, ok := s.collectionCounter[tt.args.name]
			if ok != tt.wantOk {
				t.Errorf("ClearCounterMetric() metric %s presence = %v, want %v", tt.args.name, ok, tt.wantOk)
			}
			if val != tt.want {
				t.Errorf("ClearCounterMetric() metric %s value = %v, want %v", tt.args.name, val, tt.want)
			}

			if tt.name == "clear non-existing counter metric" {
				if otherVal, otherOk := s.collectionCounter["other_counter"]; !otherOk || otherVal != tt.want {
					t.Errorf("ClearCounterMetric() affected other metric 'other_counter' = %v, want %v", otherVal, tt.want)
				}
			}
		})
	}
}
func TestClearAll(t *testing.T) {
	tests := []struct {
		name     string
		preSetup func(s *MemStorage)
	}{
		{
			name: "clear non-empty storage",
			preSetup: func(s *MemStorage) {
				s.collectionGauge["gauge1"] = 10.5
				s.collectionGauge["gauge2"] = 20.0
				s.collectionCounter["counter1"] = 5
				s.collectionCounter["counter2"] = 15
			},
		},
		{
			name: "clear empty storage",
			preSetup: func(s *MemStorage) {
			},
		},
		{
			name: "clear storage with only gauge metrics",
			preSetup: func(s *MemStorage) {
				s.collectionGauge["gauge1"] = 10.5
				s.collectionGauge["gauge2"] = 20.0
			},
		},
		{
			name: "clear storage with only counter metrics",
			preSetup: func(s *MemStorage) {
				s.collectionCounter["counter1"] = 5
				s.collectionCounter["counter2"] = 15
			},
		},
	}

	s, err := New(config.Parameters{})
	if err != nil {
		logger.LogError(err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ClearAll() // Clear before setup to ensure clean state

			if tt.preSetup != nil {
				tt.preSetup(s)
			}

			s.ClearAll()

			if len(s.collectionGauge) != 0 {
				t.Errorf("ClearAll() gauge collection not empty, length = %d", len(s.collectionGauge))
			}
			if len(s.collectionCounter) != 0 {
				t.Errorf("ClearAll() counter collection not empty, length = %d", len(s.collectionCounter))
			}
		})
	}
}
