package storage

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMetricsFromFile(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		// Создаем временный файл
		tmpFile, err := os.CreateTemp("", "test_empty")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		loaded, err := loadMetricsFromFile(tmpFile.Name())
		require.NoError(t, err)
		assert.Empty(t, loaded)
	})

	t.Run("valid metrics", func(t *testing.T) {
		testMetrics := []*metrics.Metrics{
			{
				ID:    "test1",
				MType: "gauge",
				Value:  utils.FloatToPointerFloat(1.23),
				Delta: utils.IntToPointerInt(0),
			},
			{
				ID:    "test2",
				MType: "counter",
				Delta: utils.IntToPointerInt(42),
				Value: utils.FloatToPointerFloat(0),
			},
		}

		tmpFile, err := os.CreateTemp("", "test_valid")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		data, err := json.Marshal(testMetrics)
		require.NoError(t, err)
		err = os.WriteFile(tmpFile.Name(), data, 0666)
		require.NoError(t, err)

		loaded, err := loadMetricsFromFile(tmpFile.Name())
		require.NoError(t, err)
		require.Len(t, loaded, 2)
		assert.Equal(t, testMetrics[0].ID, loaded[0].ID)
		assert.Equal(t, testMetrics[0].MType, loaded[0].MType)
		assert.Equal(t, *testMetrics[0].Value, *loaded[0].Value)
		assert.Equal(t, *testMetrics[1].Delta, *loaded[1].Delta)
	})

	t.Run("invalid file content", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_invalid")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Записываем мусор в файл
		err = os.WriteFile(tmpFile.Name(), []byte("invalid json"), 0666)
		require.NoError(t, err)

		_, err = loadMetricsFromFile(tmpFile.Name())
		assert.Error(t, err)
	})


}

func TestSaveMetricsToFile(t *testing.T) {
	t.Run("save and load roundtrip", func(t *testing.T) {
		testMetrics := []metrics.Metrics{
			{
				ID:    "roundtrip1",
				MType: "gauge",
				Value: utils.FloatToPointerFloat(3.14),
			},
			{
				ID:    "roundtrip2",
				MType: "counter",
				Delta: utils.IntToPointerInt(100),
			},
		}

		tmpFile, err := os.CreateTemp("", "test_save")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		err = saveMetricsToFile(tmpFile.Name(), &testMetrics)
		require.NoError(t, err)

		loaded, err := loadMetricsFromFile(tmpFile.Name())
		require.NoError(t, err)
		require.Len(t, loaded, 2)
		assert.Equal(t, testMetrics[0].ID, loaded[0].ID)
		assert.Equal(t, *testMetrics[0].Value, *loaded[0].Value)
		assert.Equal(t, *testMetrics[1].Delta, *loaded[1].Delta)
	})

	t.Run("empty metrics", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_empty_save")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		emptyMetrics := make([]metrics.Metrics, 0)
		err = saveMetricsToFile(tmpFile.Name(), &emptyMetrics)
		require.NoError(t, err)

		info, err := os.Stat(tmpFile.Name())
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))

		loaded, err := loadMetricsFromFile(tmpFile.Name())
		require.NoError(t, err)
		assert.Empty(t, loaded)
	})

	t.Run("invalid path", func(t *testing.T) {
		invalidPath := "/invalid/path/to/file.json"
		testMetrics := []metrics.Metrics{
			{
				ID:    "test",
				MType: "gauge",
				Value: utils.FloatToPointerFloat(1.0),
			},
		}
		err := saveMetricsToFile(invalidPath, &testMetrics)
		assert.Error(t, err)
	})
}

func TestWithSyncLocalStorage(t *testing.T) {

	t.Run("middleware calls next handler", func(t *testing.T) {
		called := false
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		origInterval := saveInterval
		origPath := localStoragePath
		defer func() {
			saveInterval = origInterval
			localStoragePath = origPath
		}()

		saveInterval = 0
		tmpFile, err := os.CreateTemp("", "test_middleware")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		localStoragePath = tmpFile.Name()

		wrapped := WithSyncLocalStorage(mockHandler)

		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		wrapped.ServeHTTP(nil, req)

		assert.True(t, called)

		_, err = os.Stat(tmpFile.Name())
		assert.NoError(t, err)
	})
}
