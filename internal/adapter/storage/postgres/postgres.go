package postgres

import (
	"github.com/jmoiron/sqlx"
)

type postgres struct {
	db           *sqlx.DB
	initialWords []string
}

func New(db *sqlx.DB, initialWords []string) *postgres {
	return &postgres{
		db:           db,
		initialWords: initialWords,
	}
}
