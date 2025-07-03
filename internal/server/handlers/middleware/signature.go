package middleware

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
)

type hashResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
}

func (r *hashResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, err
	}
	hash, err := signature.Get(b)
	if err == nil {
		encodedHash := base64.StdEncoding.EncodeToString(hash)
		r.Header().Set("HashSHA256", encodedHash)
	}
	return size, err
}

func SignatureHandle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		headerValues := r.Header.Get("HashSHA256")

		if signature.GetKey() == "" || headerValues == "" {
			next.ServeHTTP(res, r)
			return
		}
		decodedHeader, err := base64.StdEncoding.DecodeString(headerValues)
		if err != nil {
			http.Error(res, "invalid base64 encoding", http.StatusBadRequest)
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := signature.Check(decodedHeader, bodyBytes); err != nil {
			logger.LogError(err)
			http.Error(res, "", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		w := hashResponseWriter{
			ResponseWriter: res, // встраиваем оригинальный http.ResponseWriter
		}
		next.ServeHTTP(&w, r)

	})
}
