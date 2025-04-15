package client

import (
	"bytes"
	"compress/gzip"
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
			logger.LogInfo(err)

			return err
	}
		path := fmt.Sprintf("http://%s/update/", address)
		// реализует io.Writer и io.Reader
		var compressedBody bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBody)
		_, err = gzipWriter.Write(body)
		if err != nil {
			logger.LogInfo(err)
			return err
		}
		gzipWriter.Close()

		req, err := http.NewRequest("POST", path, &compressedBody)

		if err != nil {
			logger.LogInfo(err)
			return nil
		}  
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			logger.LogInfo(err)
			return err
		}
		resp.Body.Close()
	}
	return nil
}
