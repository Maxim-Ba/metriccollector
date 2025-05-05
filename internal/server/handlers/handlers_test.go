package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		method string
		body metrics.Metrics
		want   want
	}{
		{
			name:   "If Method not POST must be error",
			method: http.MethodGet,
			body: metrics.Metrics{
				ID: "name" ,
				Value:utils.IntToPointerFloat(234) ,
				MType: constants.Gauge,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "No name metrics",
			body: metrics.Metrics{
				ID: "" ,
				Value:utils.IntToPointerFloat(234) ,
				MType: constants.Gauge,
			},
			method: http.MethodPost,
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "Wrong metric type",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID: "name" ,
				Value:utils.IntToPointerFloat(234) ,
				MType: "",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},

		{
			name:   "Ok counter",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID: "name" ,
				MType: constants.Counter,
				Delta:utils.IntToPointerInt(345)  ,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "Ok gauge",
			body: metrics.Metrics{
				ID: "name" ,
				Value:utils.IntToPointerFloat(234) ,
				MType: constants.Gauge,
			},
			method: http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		},
	}
	handler := http.HandlerFunc(UpdateHandler)
	srv := httptest.NewServer(handler)
	client := &http.Client{}
	for _, test := range tests {
		path := "/update/"
		body ,err:= json.Marshal(test.body)
		if err != nil {
			logger.LogError("error in convert bpdy to JSON")
			return 
		}
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, srv.URL+path, bytes.NewReader(body))
			assert.NoError(t,err)
			res, err := client.Do(request)  
			assert.NoError(t,err)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer func() {
				if err := res.Body.Close(); err != nil {
					logger.LogError(err)
				}
			}()
		})
	}
}

func Test_getAllHandler(t *testing.T) {
	type want struct {
		code int
		contentType string

	}
	tests := []struct {
		name string
		path   string
		method string
		want   want
	}{
		{
			name: "Correct request, Ok",
			path: "/",
			method: http.MethodGet,
			want: want{
				code: http.StatusOK,
				contentType: "text/html",
			},
		},
		{
			name: "Wrong method",
			path: "/",
			method: http.MethodPost,
			want: want{
				code: http.StatusMethodNotAllowed,
				contentType: "text/html",
			},
		},
	}
	handler := http.HandlerFunc(GetAllHandler)
	srv := httptest.NewServer(handler)
	client := &http.Client{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, srv.URL+ test.path, nil)
			assert.NoError(t,err)

			res, err := client.Do(request)  
			assert.NoError(t,err)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer func() {
				if err := res.Body.Close(); err != nil {
					logger.LogError(err)
				}
			}()
		})
	}
}
