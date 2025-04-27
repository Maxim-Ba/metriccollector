package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func TestHTTPClient_SendMetrics(t *testing.T) {
	type args struct {
		metrics []*metrics.Metrics
	}
	tests := []struct {
		name        string
		client      *HTTPClient
		args        args
		wantErr     bool
		wantHeaders map[string]string
	}{
		{
			name: "Valid Content-Type",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{
					{},
				},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Valid Content-Encoding",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{{}},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Content-Encoding": "gzip",
			},
		},
		{
			name: "Valid Accept-Encoding",
			client: &HTTPClient{
				httpClient: &http.Client{},
			},
			args: args{
				metrics: []*metrics.Metrics{{}},
			},
			wantErr: false,
			wantHeaders: map[string]string{
				"Accept-Encoding": "gzip",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, value := range tt.wantHeaders {
					if r.Header.Get(key) != value {
						t.Errorf("Expected header %s to be %s, got %s", key, value, r.Header.Get(key))
					}
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			address = server.URL[7:] // Remove "http://"
			if err := tt.client.SendMetrics(tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.SendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
