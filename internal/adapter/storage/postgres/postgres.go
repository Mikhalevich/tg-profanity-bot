package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func New(connection string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", connection)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) Close() error {
	if err := p.db.Close(); err != nil {
		return fmt.Errorf("close connection: %w", err)
	}

	return nil
}
