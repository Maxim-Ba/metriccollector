package postgres

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	dbInstance       *sql.DB
	saveMetricsMutex sync.Mutex // Мьютекс для синхронизации
)

// Custom error type for unique violation
var ErrUniqueViolation = errors.New("unique violation")

func New(connectionParams string) (*sql.DB, error) {
	logger.LogInfo("postgres New")
	database, err := sql.Open("pgx", connectionParams)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	err = utils.RetryWrapper(func() error {
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR(255) PRIMARY KEY,
			type VARCHAR(255) NOT NULL,
			value DOUBLE PRECISION,
			delta BIGINT,
			CONSTRAINT chk_value_delta CHECK ((value IS NULL) OR (delta IS NULL))
		)`)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrUniqueViolation
			}
		}
		return err
	}, 3, []error{ErrUniqueViolation})

	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	dbInstance = database
	return database, nil
}

func LoadMetricsFromDB() ([]*metrics.Metrics, error) {
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
	}, 3, []error{sql.ErrConnDone, ErrUniqueViolation})

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return metricsList, nil
}

func SaveMetricsToDB(metricsList *[]metrics.Metrics) error {
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
				tx.Rollback()
				return err
			}
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}, 3, []error{sql.ErrConnDone, ErrUniqueViolation})

	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.LogInfo("bd updated")

	return nil
}
