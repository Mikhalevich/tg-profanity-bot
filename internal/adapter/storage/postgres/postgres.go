package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func New(db *sql.DB, driverName string) *Postgres {
	return &Postgres{
		db: sqlx.NewDb(db, driverName),
	}
}
