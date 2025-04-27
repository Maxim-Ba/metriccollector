package storage

import (
	"context"
	"database/sql"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage/database/postgres"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type MemStorage struct {
	collectionGauge   map[string]float64
	collectionCounter map[string]int64
}

var StorageInstance = MemStorage{
	collectionGauge:   map[string]float64{},
	collectionCounter: map[string]int64{},
}

var db *sql.DB

func New(cfg config.Parameters) (*MemStorage, error) {
	logger.LogInfo("storage New")
	initStoreValues := []*metrics.Metrics{}
	saveInterval = cfg.StoreIntervalSecond
	localStoragePath = cfg.StoragePath
	databaseDSN = cfg.DatabaseDSN
	var err error
	if cfg.Restore {
		if cfg.DatabaseDSN != "" {
			db, err = postgres.New(cfg.DatabaseDSN)
			if err != nil {
				logger.LogError(err)
				return nil, err
			}
			initStoreValues, err = postgres.LoadMetricsFromDB()
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
func Close() {
	err := db.Close()
	if err != nil {
		logger.LogError(err)

	}
}

func (s MemStorage) SaveMetric(m *metrics.Metrics) error {

	if m.MType == "gauge" {
		StorageInstance.collectionGauge[m.ID] = *m.Value
	}
	if m.MType == "counter" {
		metricValue := int64(*m.Delta)
		if val, ok := StorageInstance.collectionCounter[m.ID]; ok {
			StorageInstance.collectionCounter[m.ID] = val + metricValue
		} else {
			StorageInstance.collectionCounter[m.ID] = metricValue
		}
	}
	return nil
}

//

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
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: metric, Value: utils.FloatToPointerFloat(value)})
		}
		for metric, value := range StorageInstance.collectionCounter {
			metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: metric, Delta: utils.FloatToPointerInt(value)})
		}
		return &metricsSlice, nil
	}

	//Get choosen metrics
	for _, metric := range *metricsParams {
		if metric.MetricType == "gauge" {
			if value, ok := StorageInstance.collectionGauge[metric.MetricsName]; ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "gauge", ID: metric.MetricsName, Value: utils.FloatToPointerFloat(value)})
			}
		} else if metric.MetricType == "counter" {
			if value, ok := StorageInstance.collectionCounter[metric.MetricsName]; ok {
				metricsSlice = append(metricsSlice, metrics.Metrics{MType: "counter", ID: metric.MetricsName, Delta: utils.FloatToPointerInt(value)})
			}
		}
	}
	if len(metricsSlice) == 0 {
		return nil, ErrUnknownMetricName
	}
	return &metricsSlice, nil
}

func Ping(ctx context.Context) error {
	err := db.PingContext(ctx)
	return err
}
