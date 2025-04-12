package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
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

func (c *HTTPClient) SendMetrics(metrics []*metrics.Metrics) error {
	fmt.Println("send request")

	for _, metric := range metrics {
		body, err := json.Marshal(*metric)
		if err != nil {
			logger.LogInfo("agent 0------------")
			logger.LogInfo(err)

			return err
	}
		path := fmt.Sprintf("http://%s/update/", address)
		resp,err:=c.httpClient.Post(path, "application/json", bytes.NewReader(body))
		if err != nil {
			logger.LogInfo("agent 1------------")
			logger.LogInfo(err)
			return nil
		}  
		resp.Body.Close()
	}
	return nil
}
