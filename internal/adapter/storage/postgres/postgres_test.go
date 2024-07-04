package postgres

import (
	"database/sql"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

type PostgresSuit struct {
	*suite.Suite
	dbCleanup func() error
	p         *Postgres
}

func TestPostgresSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &PostgresSuit{
		Suite: new(suite.Suite),
	})
}

func (s *PostgresSuit) SetupSuite() {
	db, cleanup, err := connectToDatabase()
	if err != nil {
		s.FailNow("could not connect to database", err)
	}

	if err := migrationsUp(db, "../../../../script/db/migrations"); err != nil {
		s.FailNow("could not exec migrations", err)
	}

	s.dbCleanup = cleanup
	s.p = New(db, "pgx")
}

func (s *PostgresSuit) TearDownSuite() {
	if err := s.dbCleanup(); err != nil {
		s.FailNow("could not db cleanup: %v", err)
	}
}

func connectToDatabase() (*sql.DB, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("construct pool: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return nil, nil, fmt.Errorf("connect to Docker: %w", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.3-alpine3.20",
		Env: []string{
			"POSTGRES_DB=bot",
			"POSTGRES_USER=bot",
			"POSTGRES_PASSWORD=bot",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	var db *sql.DB

	if err := pool.Retry(func() error {
		db, err = sql.Open("pgx",
			fmt.Sprintf("host=localhost port=%s user=bot password=bot dbname=bot sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return fmt.Errorf("sql open: %w", err)
		}

		if err := db.Ping(); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, func() error {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("purge resource: %w", err)
		}

		return nil
	}, nil
}

func migrationsUp(db *sql.DB, pathToMigrations string) error {
	_, filename, _, _ := runtime.Caller(0)
	migrationDir, err := filepath.Abs(filepath.Join(path.Dir(filename), pathToMigrations))
	if err != nil {
		return fmt.Errorf("making migrations dir: %w", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: migrationDir,
	}

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("exec migrations: %w", err)
	}

	return nil
}
