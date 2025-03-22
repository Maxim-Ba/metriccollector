package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		method string
		want   want
	}{
		{
			name:   "If Method not POST must be error",
			path:   `gauge/name/234.34`,
			method: http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "No name metrics",
			path:   `gauge//234.34`,
			method: http.MethodPost,
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "Wrong metric type",
			path:   `g/name/234.34`,
			method: http.MethodPost,

			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong metric value",
			path:   `gauge/name/string`,
			method: http.MethodPost,

			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Ok counter",
			path:   `counter/name/345`,
			method: http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "Ok gauge",
			path:   `gauge/name/345`,
			method: http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		},
	}
	handler := http.HandlerFunc(updateHandler)
	srv := httptest.NewServer(handler)
	client := &http.Client{}
	for _, test := range tests {
		path := "/update/" + test.path
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, srv.URL+path, nil)
			res, err := client.Do(request)  
			assert.NoError(t,err)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
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
	handler := http.HandlerFunc(getAllHandler)
	srv := httptest.NewServer(handler)
	client := &http.Client{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, srv.URL+ test.path, nil)
			res, err := client.Do(request)  
			assert.NoError(t,err)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
