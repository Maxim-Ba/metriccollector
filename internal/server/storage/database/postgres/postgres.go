package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(connectionParams string) (*sql.DB, error) {

	database, err := sql.Open("pgx", connectionParams)
	if err != nil {
		return nil, err
	}
	return database, nil
}
