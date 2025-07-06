package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"slices"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
func GzipHandle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {

		r, err := decodeGzip(r)
		if err != nil {
			logger.LogError(err)
			res.WriteHeader(http.StatusMethodNotAllowed)
			_, err := res.Write([]byte(""))
			if err != nil {
				return
			}
			return
		}
		// проверяем, что клиент поддерживает gzip-сжатие
		headerValues := r.Header.Values("Accept-Encoding")

		if !slices.Contains(headerValues, "gzip") {
			next.ServeHTTP(res, r)
			return
		}
		res.Header().Set("Content-Encoding", "gzip")

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(res, gzip.BestSpeed)
		if err != nil {
			logger.LogError(err)
			return
		}
		defer func() {
			if err := gz.Close(); err != nil {
				logger.LogError(err)
			}
		}()

		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: res, Writer: gz}, r)
	})
}

func decodeGzip(r *http.Request) (*http.Request, error) {
	headerValues := r.Header.Values("Content-Encoding")

	if !slices.Contains(headerValues, "gzip") {
		return r, nil
	}
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		return r, ErrWrongBodyEncoding
	}
	defer func() {
		if err := gz.Close(); err != nil {
			logger.LogError(err)
		}
	}()

	body, err := io.ReadAll(gz)
	if err != nil {
		return r, ErrWrongBodyEncoding

	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return r, nil
}
