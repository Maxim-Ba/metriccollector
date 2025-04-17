package storage

import (
	"reflect"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		metricsNames *[]*metrics.MetricDTOParams
	}
	tests := []struct {
		name    string
		args    args
		want    *[]metrics.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	s,_ := New(config.Parameters{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetMetrics(tt.args.metricsNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

