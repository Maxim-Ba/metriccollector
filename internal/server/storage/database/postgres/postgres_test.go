package postgres

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMigrate struct {
	upError error
}

func (m *MockMigrate) Up() error {
	return m.upError
}

func (m *MockMigrate) Close() (error, error) {
	return nil, nil
}

func (m *MockMigrate) Steps(int) error {
	return nil
}

func (m *MockMigrate) Down() error {
	return nil
}

func TestNew_MigrationsUpError(t *testing.T) {
	tempDir := t.TempDir()
	migrationsDir := filepath.Join(tempDir, "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err)

	originalSQLOpen := sqlOpen
	defer func() { sqlOpen = originalSQLOpen }()
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		db, _, err := sqlmock.New()
		return db, err
	}

	originalMigrateNew := migrateNew
	defer func() { migrateNew = originalMigrateNew }()
	migrateNew = func(sourceURL, databaseURL string) (*migrate.Migrate, error) {
		return &migrate.Migrate{}, nil
	}

	db, err := New("postgres://user:pass@localhost:5432/db", migrationsDir)
	assert.Error(t, err)
	assert.Nil(t, db)
}

var (
	sqlOpen    = sql.Open
	migrateNew = migrate.New
)

func TestSaveMetricsToDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testMetrics := []metrics.Metrics{
		{ID: "test1", MType: "gauge", Value: utils.FloatToPointerFloat(1.23)},
		{ID: "test2", MType: "counter", Delta: utils.IntToPointerInt(42)},
	}

	t.Run("successful save", func(t *testing.T) {
		mock.ExpectBegin()

		for _, m := range testMetrics {
			mock.ExpectExec(`INSERT INTO metrics`).
				WithArgs(m.ID, m.MType, m.Value, m.Delta).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit()

		err := SaveMetricsToDB(&testMetrics, db)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transaction begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		err := SaveMetricsToDB(&testMetrics, db)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("commit error", func(t *testing.T) {
		mock.ExpectBegin()
		for _, m := range testMetrics {
			mock.ExpectExec(`INSERT INTO metrics`).
				WithArgs(m.ID, m.MType, m.Value, m.Delta).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := SaveMetricsToDB(&testMetrics, db)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLoadMetricsFromDB(t *testing.T) {
	expectedMetrics := []*metrics.Metrics{
		{ID: "test1", MType: "gauge", Value: utils.FloatToPointerFloat(1.23)},
		{ID: "test2", MType: "counter", Delta: utils.IntToPointerInt(42)},
	}

	t.Run("successful load", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
			AddRow("test1", "gauge", 1.23, nil).
			AddRow("test2", "counter", nil, 42)

		mock.ExpectQuery(`SELECT \* FROM metrics`).WillReturnRows(rows)

		metrics, err := LoadMetricsFromDB(db)
		assert.NoError(t, err)
		assert.Equal(t, expectedMetrics, metrics)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database not initialized", func(t *testing.T) {
		metrics, err := LoadMetricsFromDB(nil)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.EqualError(t, err, "database not initialized")
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
			AddRow("test1", "gauge", 1.23, "not_an_int")

		mock.ExpectQuery(`SELECT \* FROM metrics`).WillReturnRows(rows)

		metrics, err := LoadMetricsFromDB(db)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "type", "value", "delta"}).
			AddRow("test1", "gauge", 1.23, nil).
			RowError(0, sql.ErrNoRows)

		mock.ExpectQuery(`SELECT \* FROM metrics`).WillReturnRows(rows)

		metrics, err := LoadMetricsFromDB(db)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
