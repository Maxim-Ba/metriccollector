package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateHandler(t *testing.T) {
	type want struct {
		code        int
}
tests := []struct {
	name string
	path string
	method string
	want want
}{
	{
			name: "If Method not POST must be error",
			path: `gauge/name/234.34`,
			method: http.MethodGet,
			want: want{
					code:        http.StatusMethodNotAllowed,
			},
	},
	{
		name: "No name metrics",
		path: `gauge//234.34`,
		method: http.MethodPost,
		want: want{
				code:        http.StatusNotFound,
		},
},
{
	name: "Wrong metric type",
	path: `g/name/234.34`,
	method: http.MethodPost,

	want: want{
			code:        http.StatusBadRequest,
	},
},
{
	name: "Wrong metric value",
	path: `gauge/name/string`,
	method: http.MethodPost,

	want: want{
			code:        http.StatusBadRequest,
	},
},
{
	name: "Ok counter",
	path: `counter/name/345`,
	method: http.MethodPost,
	want: want{
			code:        http.StatusOK,
	},
},
{
	name: "Ok gauge",
	path: `gauge/name/345`,
	method: http.MethodPost,
	want: want{
			code:        http.StatusOK,
	},
},
}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := "/update/" + test.path
			request := httptest.NewRequest(test.method,path , nil)
			w := httptest.NewRecorder()
			updateHandler(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()	
		})
	}
}
