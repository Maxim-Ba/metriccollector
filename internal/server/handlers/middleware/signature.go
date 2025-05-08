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

func SignatureHandle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		logger.LogInfo("SignatureHandle")
		headerValues := r.Header.Get("HashSHA256")
		logger.LogInfo("signature.GetKey()", signature.GetKey())

		if signature.GetKey() == "" || r.Method == http.MethodGet || headerValues==""{
			next.ServeHTTP(res, r)
			return
		}
    decodedHeader, err := base64.StdEncoding.DecodeString(headerValues)
		logger.LogInfo("headerValues ", headerValues)
		logger.LogInfo("decodedHeader ", decodedHeader)

    if err != nil {
      res.WriteHeader(http.StatusBadRequest)
			_, err := res.Write([]byte("invalid base64 encoding"))

      if err != nil {
				return
			}
      return 
    }
		bodyBytes, _ := io.ReadAll(r.Body)
		if !signature.Check(decodedHeader, bodyBytes) {
			// logger.LogError("check hash fail")
			res.WriteHeader(http.StatusBadRequest)
			_, err := res.Write([]byte(""))
			if err != nil {
				return
			}
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		w := hashResponseWriter{
			ResponseWriter: res, // встраиваем оригинальный http.ResponseWriter
		}
		next.ServeHTTP(&w, r)

	})
}
func (r *hashResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, err
	}
	hash, err := signature.Get(b)
	if err == nil {
    encodedHash := base64.StdEncoding.EncodeToString(hash)
logger.LogInfo("--    ",encodedHash)
		r.Header().Set("HashSHA256", encodedHash)
	}
	return size, err
}
