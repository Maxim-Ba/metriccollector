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

var address string

func NewClient(initAddress string) *HTTPClient {
	address = initAddress
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPClient) SendMetrics(metrics []*metrics.MetricDTO) error {
	fmt.Println("send request")
	for _, metric := range metrics {
		path := fmt.Sprintf("http://%s/update/%s/%s/%f", address, metric.MetricType, metric.MetricName, metric.Value)
		resp,err:=c.httpClient.Post(path, "text/plain", nil)
		if err != nil {
			return nil
		}  
		resp.Body.Close()
	}
	return nil
}
