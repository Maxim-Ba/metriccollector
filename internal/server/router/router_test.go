package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddlewaresOrder(t *testing.T) {
	r := New()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Проверяем наличие заголовков, которые должны добавляться middleware
	assert.NotEmpty(t, rec.Header().Get("Content-Type"))

}

func TestChiRouter(t *testing.T) {
	r := New()

	// Проверяем, что возвращается именно chi.Mux
	assert.IsType(t, &chi.Mux{}, r)
}

func TestRoutePatterns(t *testing.T) {
	r := New()

	tests := []struct {
		method string
		path   string
		expect string
	}{
		{http.MethodGet, "/value/gauge/test_metric", "/value/{metricType}/{metricName}"},
		{http.MethodPost, "/update/counter/test_metric/10", "/update/{metricType}/{metricName}/{value}"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			_, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			// Создаем контекст маршрутизации
			rctx := chi.NewRouteContext()

			// Пытаемся найти маршрут
			matched := r.Match(rctx, tt.method, tt.path)
			require.True(t, matched, "expected to match route")

			// Получаем шаблон маршрута
			routePattern := rctx.RoutePattern()
			assert.Equal(t, tt.expect, routePattern)
		})
	}
}
