package pgbutton_test

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgbutton"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
)

type ButtonSuit struct {
	*suite.Suite

	transactor pgbutton.Transactor
	dbCleanup  func() error
	pgButton   *pgbutton.PgButton
}

func TestButtonSuit(t *testing.T) {
	t.Parallel()

	suite.Run(t, &ButtonSuit{
		Suite: new(suite.Suite),
	})
}

func (s *ButtonSuit) SetupSuite() {
	dbDriver := driver.NewPgx()

	dbConn, cleanup, err := connectToDatabase(s.T().Context(), dbDriver.Name())
	if err != nil {
		s.FailNow("could not connect to database", err)
	}

	if err := migrationsUp(dbConn, "../../../../../script/db/migrations"); err != nil {
		s.FailNow("could not exec migrations", err)
	}

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, dbDriver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		pgButton            = pgbutton.New(transactionProvider)
	)

	s.transactor = transactionProvider
	s.dbCleanup = cleanup
	s.pgButton = pgButton
}

func (s *ButtonSuit) TearDownSuite() {
	if err := s.dbCleanup(); err != nil {
		s.FailNow("could not db cleanup", err)
	}
}

func (s *ButtonSuit) TearDownTest() {
	s.cleanup()
}

func (s *ButtonSuit) TearDownSubTest() {
	s.cleanup()
}

func (s *ButtonSuit) cleanup() {
	var (
		ctx = context.Background()
	)

	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM buttons")
}

func connectToDatabase(ctx context.Context, driverName string) (*sql.DB, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("construct pool: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return nil, nil, fmt.Errorf("connect to docker: %w", err)
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

	if err != nil {
		return nil, nil, fmt.Errorf("run docker: %w", err)
	}

	var dbConn *sql.DB

	if err := pool.Retry(func() error {
		dbConn, err = sql.Open(driverName,
			fmt.Sprintf("host=localhost port=%s user=bot password=bot dbname=bot sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return fmt.Errorf("sql open: %w", err)
		}

		if err := dbConn.PingContext(ctx); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("connect to database: %w", err)
	}

	return dbConn, func() error {
		if err := dbConn.Close(); err != nil {
			return fmt.Errorf("close database connection: %w", err)
		}

		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("purge resource: %w", err)
		}

		return nil
	}, nil
}

func migrationsUp(dbConn *sql.DB, pathToMigrations string) error {
	//nolint:dogsled
	_, filename, _, _ := runtime.Caller(0)
	migrationDir, err := filepath.Abs(filepath.Join(path.Dir(filename), pathToMigrations))

	if err != nil {
		return fmt.Errorf("making migrations dir: %w", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: migrationDir,
	}

	_, err = migrate.Exec(dbConn, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("exec migrations: %w", err)
	}

	return nil
}
