package middleware

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"

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

// SignatureHandle is a middleware that handles request/response signing.
// For incoming requests, it verifies the HashSHA256 header if present.
// For responses, it calculates and sets the HashSHA256 header when
// a signing key is configured.
func SignatureHandle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		headerValues := r.Header.Get("HashSHA256")
		if signature.GetKey() == "" && signature.GetPrivKey() == nil {
			next.ServeHTTP(res, r)
			return
		}

		// Если есть ключ, но нет заголовка - тоже пропускаем
		if headerValues == "" && signature.GetPrivKey() == nil {
			next.ServeHTTP(res, r)
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(res, "failed to read request body", http.StatusBadRequest)
			return
		}

		if signature.GetPrivKey() != nil {
			bodyBytes, err = signature.Decrypt(bodyBytes)
			if err != nil {
				http.Error(res, "failed to decrypt body", http.StatusBadRequest)
				return
			}
		}

		// Проверяем подпись, если есть ключ подписи
		if signature.GetKey() != "" && headerValues != "" {
			decodedHeader, err := base64.StdEncoding.DecodeString(headerValues)
			if err != nil {
				http.Error(res, "invalid base64 encoding", http.StatusBadRequest)
				return
			}
			if err := signature.Check(decodedHeader, bodyBytes); err != nil {
				http.Error(res, "", http.StatusBadRequest)
				return
			}
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		w := hashResponseWriter{
			ResponseWriter: res,
		}
		next.ServeHTTP(&w, r)
	})
}
