package storage

import (
	"context"
	"database/sql"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage/database/postgres"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

// MemStorage represents an in-memory storage implementation for metrics.
// It maintains two separate collections for gauge and counter metrics.
type MemStorage struct {
	collectionGauge   map[string]float64
	collectionCounter map[string]int64
}

// StorageInstance is the global instance of MemStorage initialized with empty collections.
var StorageInstance = MemStorage{
	collectionGauge:   map[string]float64{},
	collectionCounter: map[string]int64{},
}

var db *sql.DB

// New initializes the storage system based on configuration parameters.
// It handles:
// - Database connection setup if DSN is provided
// - Metric restoration from file or database
// - Background saving routine
// Parameters:
//   - cfg: Configuration parameters including storage options
//
// Returns:
//   - *MemStorage: Initialized storage instance
//   - error: if initialization fails
func New(cfg config.Parameters) (*MemStorage, error) {
	logger.LogInfo("storage New")
	initStoreValues := []*metrics.Metrics{}
	saveInterval = cfg.StoreIntervalSecond
	localStoragePath = cfg.StoragePath
	databaseDSN = cfg.DatabaseDSN
	var err error
	if cfg.Restore {
		if cfg.DatabaseDSN != "" {
			db, err = postgres.New(cfg.DatabaseDSN, cfg.MigrationsPath)
			if err != nil {
				logger.LogError(err)
				return nil, err
			}
			initStoreValues, err = postgres.LoadMetricsFromDB(db)
		} else {
			initStoreValues, err = loadMetricsFromFile(localStoragePath)
		}
	}
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	for _, m := range initStoreValues {
		err := StorageInstance.SaveMetric(m)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}
	}
	go saveLoop()
	return &StorageInstance, nil
}

// Close terminates the database connection if it exists.
// Logs any errors encountered during closure.
func Close() {
	if db != nil {

		err := db.Close()
		if err != nil {
			logger.LogError(err)
		} else {
			logger.LogInfo("Database connection is already closed or was never opened")
		}
	}
}

// SaveMetric persists a single metric to memory storage.
// For gauge metrics, it overwrites the existing value.
// For counter metrics, it increments the existing value.
// Parameters:
//   - m: Metric to save
//
// Returns:
//   - error: if metric type is invalid
func (s MemStorage) SaveMetric(m *metrics.Metrics) error {

	if m.MType == constants.Gauge {
		StorageInstance.collectionGauge[m.ID] = *m.Value
	}
	if m.MType == constants.Counter {
		metricValue := int64(*m.Delta)
		if val, ok := StorageInstance.collectionCounter[m.ID]; ok {
			StorageInstance.collectionCounter[m.ID] = val + metricValue
		} else {
			StorageInstance.collectionCounter[m.ID] = metricValue
		}
	}
	return nil
}

// SaveMetrics persists multiple metrics to memory storage in batch.
// Delegates to SaveMetric for each individual metric.
// Parameters:
//   - metricsSlice: Slice of metrics to save
//
// Returns:
//   - error: if any metric fails to save
func (s MemStorage) SaveMetrics(metricsSlice *[]metrics.Metrics) error {
	for _, m := range *metricsSlice {
		err := StorageInstance.SaveMetric(&m)
		if err != nil {
			logger.LogError(err)
			return err
		}
	}

	return nil
}

// GetMetrics retrieves metrics based on provided parameters.
// Behavior:
// - With empty params: returns all metrics
// - With specific params: returns only requested metrics
// Parameters:
//   - metricsParams: Slice of metric lookup parameters
//
// Returns:
//   - *[]metrics.Metrics: Retrieved metrics
//   - error: if no metrics found (with specific params)
func (s MemStorage) GetMetrics(metricsParams *[]*metrics.MetricDTOParams) (*[]metrics.Metrics, error) {

	metricsNames := make([]string, len(*metricsParams))
	metricsTypes := make([]string, len(*metricsParams))
	for i, m := range *metricsParams {
		metricsNames[i] = m.MetricsName
		metricsTypes[i] = m.MetricType
	}
	var metricsSlice []metrics.Metrics
	// Get all metrics
	if len(metricsNames) == 0 {
		for metric, value := range StorageInstance.collectionGauge {
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: constants.Gauge, ID: metric, Value: utils.FloatToPointerFloat(value)})
		}
		for metric, value := range StorageInstance.collectionCounter {
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: constants.Counter, ID: metric, Delta: utils.FloatToPointerInt(value)})
		}
		return &metricsSlice, nil
	}

	//Get choosen metrics
	for _, metric := range *metricsParams {
		if metric.MetricType == constants.Gauge {
			if value, ok := StorageInstance.collectionGauge[metric.MetricsName]; ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: constants.Gauge, ID: metric.MetricsName, Value: utils.FloatToPointerFloat(value)})
			}
		} else if metric.MetricType == constants.Counter {
			if value, ok := StorageInstance.collectionCounter[metric.MetricsName]; ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: constants.Counter, ID: metric.MetricsName, Delta: utils.FloatToPointerInt(value)})
			}
		}
	}
	if len(metricsSlice) == 0 {
		return nil, ErrUnknownMetricName
	}
	return &metricsSlice, nil
}

// Ping verifies the database connection is alive.
// Parameters:
//   - ctx: Context for operation cancellation
//
// Returns:
//   - error: if connection check fails or no connection exists
func (s MemStorage) Ping(ctx context.Context) error {
	if db == nil {
		return ErrDatabaseConnection
	}
	return db.PingContext(ctx)
}

// ClearGaugeMetric removes a specific gauge metric from storage.
// Parameters:
//   - name: Name of the gauge metric to remove
func (s *MemStorage) ClearGaugeMetric(name string) {
	delete(s.collectionGauge, name)
}

// ClearCounterMetric removes a specific counter metric from storage.
// Parameters:
//   - name: Name of the counter metric to remove
func (s *MemStorage) ClearCounterMetric(name string) {
	delete(s.collectionCounter, name)
}

// ClearAll resets the storage by removing all metrics.
// Reinitializes both gauge and counter collections.
func (s *MemStorage) ClearAll() {
	s.collectionGauge = make(map[string]float64)
	s.collectionCounter = make(map[string]int64)
}
