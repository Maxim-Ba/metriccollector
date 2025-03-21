package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

type Client interface {
	SendMetrics(metrics []*metrics.MetricDTO) error
}

type HTTPClient struct {
	httpClient *http.Client
}

var port int
var host string

func NewClient(initHost string, initPort int) *HTTPClient {
	port = initPort
	host = initHost
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPClient) SendMetrics(metrics []*metrics.MetricDTO) error {
	fmt.Println("send request")
	for _, metric := range metrics {
		path := fmt.Sprintf("http://%s:%d/update/%s/%s/%f", host, port, metric.MetricType, metric.MetricName, metric.Value)
		c.httpClient.Post(path, "text/plain", nil)
	}
	return nil
}
