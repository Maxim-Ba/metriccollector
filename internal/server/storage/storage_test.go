package storage

import (
	"reflect"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		metricsNames *[]string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]metrics.MetricDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMetrics(tt.args.metricsNames)
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

func Test_checkCorrectMetricsName(t *testing.T) {
	type want struct {
		result  error
	}
	tests := []struct {
		name    string
		metricsNames *[]string
		want    want
	}{
		{
			name: "correct name case",
			metricsNames: &[]string{"StackInuse"},
			want: want{
				result : nil,
			},
		},
		{
			name: "correct names case",
			metricsNames: &[]string{"StackInuse", "PollCount"},
			want: want{
				result : nil,
			},
		},
		{
			name: "wrong names case",
			metricsNames: &[]string{"Stack4545Inuse1"},
			want: want{
				result : ErrUnknownMetricName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkCorrectMetricsName(tt.metricsNames); err != tt.want.result {
				t.Errorf("checkCorrectMetricsName() error = %v, wantErr %v", err, tt.want.result)
			}
		})
	}
}
