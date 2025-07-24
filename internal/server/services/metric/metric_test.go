package metric

import (
	"errors"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage реализует интерфейс Storage для тестирования
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveMetric(metric *metrics.Metrics) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *MockStorage) SaveMetrics(metrics *[]metrics.Metrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func (m *MockStorage) GetMetrics(params *[]*metrics.MetricDTOParams) (*[]metrics.Metrics, error) {
	args := m.Called(params)
	return args.Get(0).(*[]metrics.Metrics), args.Error(1)
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockStorage)
		expectedHTML  string
		expectedError error
	}{
		{
			name: "successful retrieval",
			mockSetup: func(m *MockStorage) {
				m.On("GetMetrics", mock.Anything).Return(&[]metrics.Metrics{
					{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
					{ID: "test2", MType: "counter", Delta: int64Ptr(42)},
				}, nil)
			},
			expectedHTML:  "", // Здесь должна быть ожидаемая HTML строка
			expectedError: nil,
		},
		{
			name: "storage error",
			mockSetup: func(m *MockStorage) {
				m.On("GetMetrics", mock.Anything).Return(&[]metrics.Metrics{}, errors.New("storage error"))
			},
			expectedHTML:  "",
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.mockSetup(mockStorage)

			html, err := GetAll(mockStorage)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Contains(t, html, "test1") // Проверяем, что HTML содержит ожидаемые метрики
				assert.Contains(t, html, "test2")
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name           string
		params         []*metrics.MetricDTOParams
		mockSetup      func(*MockStorage)
		expectedMetric *metrics.Metrics
		expectedError  error
	}{
		{
			name: "successful retrieval",
			params: []*metrics.MetricDTOParams{
				{MetricsName: "test1", MetricType: "gauge"},
			},
			mockSetup: func(m *MockStorage) {
				m.On("GetMetrics", mock.Anything).Return(&[]metrics.Metrics{
					{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
				}, nil)
			},
			expectedMetric: &metrics.Metrics{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
			expectedError:  nil,
		},
		{
			name: "storage error",
			params: []*metrics.MetricDTOParams{
				{MetricsName: "test1", MetricType: "gauge"},
			},
			mockSetup: func(m *MockStorage) {
				m.On("GetMetrics", mock.Anything).Return(&[]metrics.Metrics{}, errors.New("storage error"))
			},
			expectedMetric: nil,
			expectedError:  errors.New("storage error"),
		},
		{
			name: "empty result",
			params: []*metrics.MetricDTOParams{
				{MetricsName: "test1", MetricType: "gauge"},
			},
			mockSetup: func(m *MockStorage) {
				m.On("GetMetrics", mock.Anything).Return(&[]metrics.Metrics{}, nil)
			},
			expectedMetric: nil,
			expectedError:  errors.New("no metrics found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.mockSetup(mockStorage)

			metric, err := Get(mockStorage, &tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMetric, metric)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		metric        *metrics.Metrics
		mockSetup     func(*MockStorage)
		expectedError error
	}{
		{
			name:   "successful update",
			metric: &metrics.Metrics{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
			mockSetup: func(m *MockStorage) {
				m.On("SaveMetric", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "storage error",
			metric: &metrics.Metrics{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
			mockSetup: func(m *MockStorage) {
				m.On("SaveMetric", mock.Anything).Return(errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.mockSetup(mockStorage)

			err := Update(mockStorage, tt.metric)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestUpdateMany(t *testing.T) {
	tests := []struct {
		name          string
		metrics       []metrics.Metrics
		mockSetup     func(*MockStorage)
		expectedError error
	}{
		{
			name: "successful batch update",
			metrics: []metrics.Metrics{
				{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
				{ID: "test2", MType: "counter", Delta: int64Ptr(42)},
			},
			mockSetup: func(m *MockStorage) {
				m.On("SaveMetrics", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "storage error",
			metrics: []metrics.Metrics{
				{ID: "test1", MType: "gauge", Value: float64Ptr(1.23)},
			},
			mockSetup: func(m *MockStorage) {
				m.On("SaveMetrics", mock.Anything).Return(errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			tt.mockSetup(mockStorage)

			err := UpdateMany(mockStorage, &tt.metrics)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

// Вспомогательные функции для создания указателей на значения
func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
