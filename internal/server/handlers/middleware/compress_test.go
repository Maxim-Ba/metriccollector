package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipHandle_NoCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test response"))
		require.NoError(t, err)
	})

	gzipHandler := GzipHandle(handler)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, "test response", w.Body.String())
	assert.Empty(t, resp.Header.Get("Content-Encoding"))

	if err := resp.Body.Close(); err != nil {
		require.NoError(t, err)
	}
}

func TestGzipHandle_ResponseCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test response"))
		require.NoError(t, err)
	})

	gzipHandler := GzipHandle(handler)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gzReader, err := gzip.NewReader(w.Body)
	require.NoError(t, err)
	defer func() {
		err = gzReader.Close()
		require.NoError(t, err)
	}()

	body, err := io.ReadAll(gzReader)
	require.NoError(t, err)
	assert.Equal(t, "test response", string(body))
	if err := resp.Body.Close(); err != nil {
		require.NoError(t, err)
	}
}

func TestGzipHandle_RequestDecompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		_, err = w.Write(body)
		require.NoError(t, err)
	})

	gzipHandler := GzipHandle(handler)

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, err := gzWriter.Write([]byte("test request"))
	require.NoError(t, err)
	require.NoError(t, gzWriter.Close())

	req := httptest.NewRequest("POST", "http://example.com", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gzReader, err := gzip.NewReader(w.Body)
	require.NoError(t, err)
	defer func() {
		err = gzReader.Close()
		require.NoError(t, err)
	}()
	body, err := io.ReadAll(gzReader)
	require.NoError(t, err)
	assert.Equal(t, "test request", string(body))
	if err := resp.Body.Close(); err != nil {
		require.NoError(t, err)
	}
}

func TestGzipHandle_InvalidGzipRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("should not be called"))
		require.NoError(t, err)
	})

	gzipHandler := GzipHandle(handler)

	req := httptest.NewRequest("POST", "http://example.com", strings.NewReader("invalid gzip data"))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	gzipHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	if err := resp.Body.Close(); err != nil {
		require.NoError(t, err)
	}
}

func TestGzipWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	defer func() {
		err := gzWriter.Close()
		require.NoError(t, err)
	}()
	recorder := httptest.NewRecorder()
	writer := gzipWriter{
		ResponseWriter: recorder,
		Writer:         gzWriter,
	}

	n, err := writer.Write([]byte("test data"))
	require.NoError(t, err)
	assert.Equal(t, len("test data"), n)

	require.NoError(t, gzWriter.Close())

	gzReader, err := gzip.NewReader(&buf)
	require.NoError(t, err)
	defer func() {
		err = gzReader.Close()
		require.NoError(t, err)
	}()
	body, err := io.ReadAll(gzReader)
	require.NoError(t, err)
	assert.Equal(t, "test data", string(body))
	
}
