package client

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
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
	logger.LogInfo("send request")

	for _, metric := range metrics {
		body, err := json.Marshal(*metric)
		if err != nil {
			logger.LogError(err)
			return err
		}
		path := fmt.Sprintf("http://%s/update/", address)
		// реализует io.Writer и io.Reader
		var compressedBody bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBody)
		_, err = gzipWriter.Write(body)
		if err != nil {
			logger.LogError(err)
			return err
		}
		err = gzipWriter.Close()
		if err != nil {
			logger.LogError(err)
		}
		req, err := http.NewRequest("POST", path, &compressedBody)
		if err != nil {
			logger.LogError(err)
			return nil
		}
		if signature.GetKey() != "" {
			hash, err := signature.Get(body)
			if err != nil {
				logger.LogError(err)
				return err
			}

			encodedHash := base64.StdEncoding.EncodeToString(hash)
			req.Header.Set("HashSHA256", encodedHash)
		}
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			logger.LogError(err)
			return err
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return ErrServerInternalError
		}
		if resp.StatusCode == http.StatusRequestTimeout {
			return ErrRequestTimeout
		}
		err = resp.Body.Close()
		if err != nil {
			logger.LogError(err)
			return err
		}

	}
	return nil
}

func (c *HTTPClient) SendMetricsWithBatch(metrics []*metrics.Metrics) error {
	logger.LogInfo("send request with batch")

	body, err := json.Marshal(metrics)
	if err != nil {
		logger.LogError(err)
		return err
	}
	path := fmt.Sprintf("http://%s/updates/", address)
	// реализует io.Writer и io.Reader
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	_, err = gzipWriter.Write(body)
	if err != nil {
		logger.LogError(err)
		return err
	}
	err = gzipWriter.Close()
	if err != nil {
		logger.LogError(err)
	}
	req, err := http.NewRequest("POST", path, &compressedBody)

	if err != nil {
		logger.LogError(err)
		return nil
	}
	if signature.GetKey() != "" {
		hash, err := signature.Get(body)
		if err != nil {
			logger.LogError(err)
			return err
		}
		encodedHash := base64.StdEncoding.EncodeToString(hash)
		req.Header.Set("HashSHA256", encodedHash)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.LogError(err)
		return err
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerInternalError
	}
	if resp.StatusCode == http.StatusRequestTimeout {
		return ErrRequestTimeout
	}
	err = resp.Body.Close()
	if err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}
