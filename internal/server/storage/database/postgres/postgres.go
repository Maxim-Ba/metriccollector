package postgres

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	saveMetricsMutex sync.Mutex // Мьютекс для синхронизации
)

var ErrConnectionException = errors.New("connection exception")
var ErrConnectionFailure = errors.New("connection failure")
var ErrConnectionClosed = errors.New("connection closed")

func New(connectionParams string, migrationsPath string) (*sql.DB, error) {
	logger.LogInfo("postgres New")
	database, err := sql.Open("pgx", connectionParams)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	// Применение миграций
	m, err := migrate.New(
		"file://"+migrationsPath,
		connectionParams,
	)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.LogError(err)
		return nil, err
	}

	return database, nil
}

func LoadMetricsFromDB(dbInstance *sql.DB) ([]*metrics.Metrics, error) {
	logger.LogInfo("LoadMetricsFromDB")

	if dbInstance == nil {
		err := errors.New("database not initialized")
		logger.LogError(err)
		return nil, err
	}

	var metricsList []*metrics.Metrics
	err := utils.RetryWrapper(func() error {
		rows, err := dbInstance.Query(`SELECT * FROM metrics`)
		if err != nil {
			return err
		}
		defer func() {
			err := rows.Close()
			if err != nil {
				logger.LogError(err)
			}
		}()

		for rows.Next() {
			var m metrics.Metrics
			if err := rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta); err != nil {
				return err
			}
			metricsList = append(metricsList, &m)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		return nil
	},  []error{sql.ErrConnDone})

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return metricsList, nil
}

func SaveMetricsToDB(metricsList *[]metrics.Metrics, dbInstance *sql.DB) error {
	saveMetricsMutex.Lock()         // Заблокировать перед выполнением
	defer saveMetricsMutex.Unlock() // Разблокировать в конце

	err := utils.RetryWrapper(func() error {
		tx, err := dbInstance.Begin()
		if err != nil {
			return err
		}
		for _, m := range *metricsList {
			// все изменения записываются в транзакцию
			_, err := dbInstance.Exec(`INSERT INTO metrics (id, type, value, delta) 
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (id) DO UPDATE 
				SET type = $2, value = $3, delta = $4`,
				m.ID, m.MType, m.Value, m.Delta)
			if err != nil {
				logger.LogError(err)
				err := tx.Rollback()
				if err != nil {
					logger.LogError(err)
				}
				return err
			}
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	},  []error{sql.ErrConnDone})

	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.LogInfo("bd updated")

	return nil
}
