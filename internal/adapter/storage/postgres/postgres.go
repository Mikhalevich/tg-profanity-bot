package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
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
