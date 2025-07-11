package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	metricsService "github.com/Maxim-Ba/metriccollector/internal/server/services/metric"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
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
		body   metrics.Metrics
		want   want
	}{
		{
			name:   "If Method not POST must be error",
			method: http.MethodGet,
			body: metrics.Metrics{
				ID:    "name",
				Value: utils.IntToPointerFloat(234),
				MType: constants.Gauge,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "No name metrics",
			body: metrics.Metrics{
				ID:    "",
				Value: utils.IntToPointerFloat(234),
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
				ID:    "name",
				Value: utils.IntToPointerFloat(234),
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
				ID:    "name",
				MType: constants.Counter,
				Delta: utils.IntToPointerInt(345),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Ok gauge",
			body: metrics.Metrics{
				ID:    "name",
				Value: utils.IntToPointerFloat(234),
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
		body, err := json.Marshal(test.body)
		if err != nil {
			logger.LogError("error in convert bpdy to JSON")
			return
		}
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, srv.URL+path, bytes.NewReader(body))
			assert.NoError(t, err)
			res, err := client.Do(request)
			assert.NoError(t, err)
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
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		path   string
		method string
		want   want
	}{
		{
			name:   "Correct request, Ok",
			path:   "/",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/html",
			},
		},
		{
			name:   "Wrong method",
			path:   "/",
			method: http.MethodPost,
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "text/html",
			},
		},
	}
	handler := http.HandlerFunc(GetAllHandler)
	srv := httptest.NewServer(handler)
	client := &http.Client{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, srv.URL+test.path, nil)
			assert.NoError(t, err)

			res, err := client.Do(request)
			assert.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer func() {
				if err := res.Body.Close(); err != nil {
					logger.LogError(err)
				}
			}()
		})
	}
}

func TestGetOneHandlerByParams(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}

	tests := []struct {
		name       string
		method     string
		path       string
		want       want
		prepMetric *metrics.Metrics
	}{
		{
			name:   "Wrong method - not GET",
			method: http.MethodPost,
			path:   "/value/gauge/testGauge",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "Not enough parameters in path - only type",
			method: http.MethodGet,
			path:   "/value/gauge/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "",
			},
		},
		{
			name:   "Wrong metric type",
			method: http.MethodGet,
			path:   "/value/unknown/test",
			want: want{
				code:        http.StatusNotFound,
				contentType: "",
			},
		},
		{
			name:   "Get existing gauge metric",
			method: http.MethodGet,
			path:   "/value/gauge/testGauge",
			prepMetric: &metrics.Metrics{
				ID:    "testGauge",
				MType: constants.Gauge,
				Value: utils.FloatToPointerFloat(123.45),
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				body:        "123.45",
			},
		},
		{
			name:   "Get existing counter metric",
			method: http.MethodGet,
			path:   "/value/counter/testCounter",
			prepMetric: &metrics.Metrics{
				ID:    "testCounter",
				MType: constants.Counter,
				Delta: utils.FloatToPointerInt(42),
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				body:        "42",
			},
		},
		{
			name:   "Get non-existent metric",
			method: http.MethodGet,
			path:   "/value/gauge/nonExistent",
			want: want{
				code:        http.StatusNotFound,
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.prepMetric != nil {
				err := metricsService.Update(storage.StorageInstance, test.prepMetric)
				if err != nil {
					t.Fatalf("Failed to prepare test metric: %v", err)
				}
			}

			handler := http.HandlerFunc(GetOneHandlerByParams)
			srv := httptest.NewServer(handler)
			defer srv.Close()

			client := &http.Client{}
			request, err := http.NewRequest(test.method, srv.URL+test.path, nil)
			assert.NoError(t, err)

			res, err := client.Do(request)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.contentType != "" {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}

			if test.want.body != "" {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, test.want.body, buf.String())
			}
			if test.prepMetric != nil {
				if test.prepMetric.MType == constants.Gauge {
					storage.StorageInstance.ClearGaugeMetric(test.prepMetric.ID)
				} else {
					storage.StorageInstance.ClearCounterMetric(test.prepMetric.ID)
				}
			}
		})
	}
}

func TestGetOneHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}

	tests := []struct {
		name       string
		method     string
		body       metrics.Metrics
		want       want
		prepMetric *metrics.Metrics
	}{
		{
			name:   "Wrong method - not POST",
			method: http.MethodGet,
			body: metrics.Metrics{
				ID:    "testGauge",
				MType: constants.Gauge,
			},
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "No metric name",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID:    "",
				MType: constants.Gauge,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:   "Wrong metric type",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID:    "test",
				MType: "unknown",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:   "Get existing gauge metric",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID:    "testGauge",
				MType: constants.Gauge,
			},
			prepMetric: &metrics.Metrics{
				ID:    "testGauge",
				MType: constants.Gauge,
				Value: utils.FloatToPointerFloat(123.45),
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"testGauge","type":"gauge","value":123.45}`,
			},
		},
		{
			name:   "Get existing counter metric",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID:    "testCounter",
				MType: constants.Counter,
			},
			prepMetric: &metrics.Metrics{
				ID:    "testCounter",
				MType: constants.Counter,
				Delta: utils.FloatToPointerInt(42),
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"testCounter","type":"counter","delta":42}`,
			},
		},
		{
			name:   "Get non-existent metric",
			method: http.MethodPost,
			body: metrics.Metrics{
				ID:    "nonExistent",
				MType: constants.Gauge,
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Подготавливаем тестовые данные если нужно
			if test.prepMetric != nil {
				err := metricsService.Update(storage.StorageInstance, test.prepMetric)
				if err != nil {
					t.Fatalf("Failed to prepare test metric: %v", err)
				}
			}

			handler := http.HandlerFunc(GetOneHandler)
			srv := httptest.NewServer(handler)
			defer srv.Close()

			// Подготавливаем тело запроса
			body, err := json.Marshal(test.body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			request, err := http.NewRequest(test.method, srv.URL, bytes.NewReader(body))
			assert.NoError(t, err)

			client := &http.Client{}
			res, err := client.Do(request)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.contentType != "" {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}

			if test.want.body != "" {
				var actualBody bytes.Buffer
				_, err := actualBody.ReadFrom(res.Body)
				assert.NoError(t, err)

				// Для сравнения JSON распаковываем оба тела
				var expected, actual interface{}
				err = json.Unmarshal([]byte(test.want.body), &expected)
				assert.NoError(t, err)
				err = json.Unmarshal(actualBody.Bytes(), &actual)
				assert.NoError(t, err)

				assert.Equal(t, expected, actual)
			}

			// Очищаем хранилище после теста
			if test.prepMetric != nil {
				if test.prepMetric.MType == constants.Gauge {
					storage.StorageInstance.ClearGaugeMetric(test.prepMetric.ID)
				} else {
					storage.StorageInstance.ClearCounterMetric(test.prepMetric.ID)
				}
			}
		})
	}
}
func TestUpdateHandlerByURLParams(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	type metricCheck struct {
		name  string
		mType string
		value interface{}
	}

	tests := []struct {
		name   string
		method string
		path   string
		want   want
		checks []metricCheck // Метрики для проверки после запроса
		setup  []metricCheck // Предварительная настройка данных
	}{
		{
			name:   "Wrong method - not POST or GET",
			method: http.MethodPut,
			path:   "/update/gauge/test/123.45",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:   "Not enough parameters in path",
			method: http.MethodPost,
			path:   "/update/gauge",
			want: want{
				code:        http.StatusNotFound,
				contentType: "",
			},
		},
		{
			name:   "Wrong metric type",
			method: http.MethodPost,
			path:   "/update/unknown/test/123",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:   "Update gauge metric via POST",
			method: http.MethodPost,
			path:   "/update/gauge/testGauge/123.45",
			want: want{
				code:        http.StatusOK,
				contentType: "",
			},
			checks: []metricCheck{
				{
					name:  "testGauge",
					mType: constants.Gauge,
					value: 123.45,
				},
			},
		},
		{
			name:   "Update counter metric via POST",
			method: http.MethodPost,
			path:   "/update/counter/testCounter/42",
			want: want{
				code:        http.StatusOK,
				contentType: "",
			},
			checks: []metricCheck{
				{
					name:  "testCounter",
					mType: constants.Counter,
					value: int64(42),
				},
			},
		},
		{
			name:   "Increment counter metric",
			method: http.MethodPost,
			path:   "/update/counter/testCounter/10",
			setup: []metricCheck{
				{
					name:  "testCounter",
					mType: constants.Counter,
					value: int64(5),
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "",
			},
			checks: []metricCheck{
				{
					name:  "testCounter",
					mType: constants.Counter,
					value: int64(15),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Очищаем хранилище перед тестом
			storage.StorageInstance.ClearAll()

			// Подготавливаем начальные данные
			for _, m := range test.setup {
				switch m.mType {
				case constants.Gauge:
					val := m.value.(float64)
					metricsService.Update(storage.StorageInstance, &metrics.Metrics{
						ID:    m.name,
						MType: m.mType,
						Value: &val,
					})
				case constants.Counter:
					val := m.value.(int64)
					metricsService.Update(storage.StorageInstance, &metrics.Metrics{
						ID:    m.name,
						MType: m.mType,
						Delta: &val,
					})
				}
			}

			handler := http.HandlerFunc(UpdateHandlerByURLParams)
			srv := httptest.NewServer(handler)
			defer srv.Close()

			req, err := http.NewRequest(test.method, srv.URL+test.path, nil)
			assert.NoError(t, err)

			client := &http.Client{}
			res, err := client.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.contentType != "" {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}

			// Проверяем состояние метрик
			for _, check := range test.checks {
				metricParams := []*metrics.MetricDTOParams{
					{
						MetricsName: check.name,
						MetricType:  check.mType,
					},
				}

				result, err := metricsService.Get(storage.StorageInstance, &metricParams)
				assert.NoError(t, err)

				switch check.mType {
				case constants.Gauge:
					assert.Equal(t, check.value.(float64), *result.Value)
				case constants.Counter:
					assert.Equal(t, check.value.(int64), int64(*result.Delta))
				}
			}
		})
	}
}
func Test_metricRecord(t *testing.T) {
	tests := []struct {
		name       string
		parameters []string
		want       metrics.Metrics
		wantErr    bool
		errType    error
	}{
		{
			name:       "Not enough parameters",
			parameters: []string{"gauge", "test"},
			want:       metrics.Metrics{},
			wantErr:    true,
			errType:    ErrNoMetricName,
		},
		{
			name:       "Wrong metric type",
			parameters: []string{"invalid", "test", "123"},
			want:       metrics.Metrics{},
			wantErr:    true,
			errType:    ErrNoMetricsType,
		},
		{
			name:       "Empty metric name",
			parameters: []string{"gauge", "", "123.45"},
			want:       metrics.Metrics{},
			wantErr:    true,
			errType:    ErrNoMetricName,
		},
		{
			name:       "Invalid gauge value",
			parameters: []string{"gauge", "test", "invalid"},
			want:       metrics.Metrics{},
			wantErr:    true,
			errType:    ErrWrongValue,
		},
		{
			name:       "Invalid counter value",
			parameters: []string{"counter", "test", "invalid"},
			want:       metrics.Metrics{},
			wantErr:    true,
			errType:    ErrWrongValue,
		},
		{
			name:       "Valid gauge metric",
			parameters: []string{"gauge", "testGauge", "123.45"},
			want: metrics.Metrics{
				MType: constants.Gauge,
				ID:    "testGauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: false,
		},
		{
			name:       "Valid counter metric",
			parameters: []string{"counter", "testCounter", "42"},
			want: metrics.Metrics{
				MType: constants.Counter,
				ID:    "testCounter",
				Delta: func() *int64 { v := int64(42); return &v }(),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := metricRecord(test.parameters)

			if test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, test.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want.MType, got.MType)
				assert.Equal(t, test.want.ID, got.ID)

				if test.want.MType == constants.Gauge {
					assert.Equal(t, *test.want.Value, *got.Value)
				} else {
					assert.Equal(t, *test.want.Delta, *got.Delta)
				}
			}
		})
	}
}
func Test_checkForAllowedMethod(t *testing.T) {
	type args struct {
		req           *http.Request
		allowedMethod []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Allowed method - GET in allowed methods",
			args: args{
				req:           httptest.NewRequest(http.MethodGet, "/", nil),
				allowedMethod: []string{http.MethodGet, http.MethodPost},
			},
			wantErr: false,
		},
		{
			name: "Allowed method - POST in allowed methods",
			args: args{
				req:           httptest.NewRequest(http.MethodPost, "/", nil),
				allowedMethod: []string{http.MethodGet, http.MethodPost},
			},
			wantErr: false,
		},
		{
			name: "Not allowed method - PUT not in allowed methods",
			args: args{
				req:           httptest.NewRequest(http.MethodPut, "/", nil),
				allowedMethod: []string{http.MethodGet, http.MethodPost},
			},
			wantErr: true,
		},
		{
			name: "Not allowed method - DELETE not in allowed methods",
			args: args{
				req:           httptest.NewRequest(http.MethodDelete, "/", nil),
				allowedMethod: []string{http.MethodGet, http.MethodPost},
			},
			wantErr: true,
		},
		{
			name: "Empty allowed methods - any method should fail",
			args: args{
				req:           httptest.NewRequest(http.MethodGet, "/", nil),
				allowedMethod: []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkForAllowedMethod(tt.args.req, tt.args.allowedMethod)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "not allowed method", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestUpdatesHandler(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name    string
		method  string
		body    []metrics.Metrics
		want    want
		wantErr bool
	}{
		{
			name:   "Wrong method - not POST",
			method: http.MethodGet,
			body:   []metrics.Metrics{},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
			wantErr: true,
		},
		{
			name:   "Empty metrics slice - should be accepted",
			method: http.MethodPost,
			body:   []metrics.Metrics{},
			want: want{
				code: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name:   "Valid gauge metrics",
			method: http.MethodPost,
			body: []metrics.Metrics{
				{
					ID:    "testGauge1",
					MType: constants.Gauge,
					Value: utils.FloatToPointerFloat(123.45),
				},
			},
			want: want{
				code: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name:   "Metric with empty name",
			method: http.MethodPost,
			body: []metrics.Metrics{
				{
					ID:    "",
					MType: constants.Gauge,
					Value: utils.FloatToPointerFloat(123.45),
				},
			},
			want: want{
				code: http.StatusNotFound,
			},
			wantErr: true,
		},
		{
			name:   "Metric with wrong type",
			method: http.MethodPost,
			body: []metrics.Metrics{
				{
					ID:    "test",
					MType: "invalid",
					Value: utils.FloatToPointerFloat(123.45),
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.StorageInstance.ClearAll()

			body, err := json.Marshal(test.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(test.method, "/updates/", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			UpdatesHandler(rec, req)

			assert.Equal(t, test.want.code, rec.Code)

			if !test.wantErr && test.want.code == http.StatusOK && len(test.body) > 0 {
				for _, m := range test.body {
					metricParams := []*metrics.MetricDTOParams{
						{MetricsName: m.ID, MetricType: m.MType},
					}
					result, err := metricsService.Get(storage.StorageInstance, &metricParams)
					assert.NoError(t, err)

					switch m.MType {
					case constants.Gauge:
						assert.Equal(t, *m.Value, *result.Value)
					case constants.Counter:
						assert.Equal(t, *m.Delta, int64(*result.Delta))
					}
				}
			}
		})
	}
}

// ..............................

func BenchmarkUpdateHandler(b *testing.B) {
	handler := http.HandlerFunc(UpdateHandler)
	testMetric := metrics.Metrics{
		ID:    "testMetric",
		MType: constants.Gauge,
		Value: utils.FloatToPointerFloat(123.45),
	}
	body, _ := json.Marshal(testMetric)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkGetAllHandler(b *testing.B) {
	handler := http.HandlerFunc(GetAllHandler)

	// Предварительно заполняем хранилище
	for i := 0; i < 100; i++ {
		metric := metrics.Metrics{
			ID:    "metric" + strconv.Itoa(i),
			MType: constants.Gauge,
			Value: utils.FloatToPointerFloat(float64(i)),
		}
		_ = metricsService.Update(storage.StorageInstance, &metric)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkGetOneHandlerByParams(b *testing.B) {
	handler := http.HandlerFunc(GetOneHandlerByParams)
	// Подготовка тестовой метрики
	testMetric := metrics.Metrics{
		ID:    "benchmarkMetric",
		MType: constants.Gauge,
		Value: utils.FloatToPointerFloat(123.45),
	}
	_ = metricsService.Update(storage.StorageInstance, &testMetric)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/benchmarkMetric", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkGetOneHandler(b *testing.B) {
	handler := http.HandlerFunc(GetOneHandler)
	testMetric := metrics.Metrics{
		ID:    "testMetric",
		MType: constants.Gauge,
	}
	body, _ := json.Marshal(testMetric)

	// Подготовка тестовой метрики
	storageMetric := metrics.Metrics{
		ID:    "testMetric",
		MType: constants.Gauge,
		Value: utils.FloatToPointerFloat(123.45),
	}
	_ = metricsService.Update(storage.StorageInstance, &storageMetric)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkUpdateHandlerByURLParams(b *testing.B) {
	handler := http.HandlerFunc(UpdateHandlerByURLParams)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/update/gauge/testMetric/123.45", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkUpdatesHandler(b *testing.B) {
	handler := http.HandlerFunc(UpdatesHandler)
	metricsSlice := []metrics.Metrics{
		{
			ID:    "metric1",
			MType: constants.Gauge,
			Value: utils.FloatToPointerFloat(123.45),
		},
		{
			ID:    "metric2",
			MType: constants.Counter,
			Delta: utils.FloatToPointerInt(42),
		},
	}
	body, _ := json.Marshal(metricsSlice)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkPingDB(b *testing.B) {
	handler := http.HandlerFunc(PingDB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
