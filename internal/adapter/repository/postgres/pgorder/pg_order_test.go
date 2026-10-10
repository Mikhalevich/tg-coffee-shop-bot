package pgorder_test

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v4"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
)

type PgOrderSuit struct {
	*suite.Suite

	transactor pgorder.Transactor
	pgOrder    *pgorder.PgOrder
}

func TestPgOrderSuit(t *testing.T) {
	t.Parallel()

	suite.Run(t, &PgOrderSuit{
		Suite: new(suite.Suite),
	})
}

func (s *PgOrderSuit) SetupSuite() {
	dbDriver := driver.NewPgx()

	dbConn, err := connectToDatabase(s.T(), dbDriver.Name())
	if err != nil {
		s.FailNow("could not connect to database", err)
	}

	if err := migrationsUp(dbConn, "../../../../script/db/migrations"); err != nil {
		s.FailNow("could not exec migrations", err)
	}

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, dbDriver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		pgOrder             = pgorder.New(dbDriver, transactionProvider)
	)

	s.transactor = transactionProvider
	s.pgOrder = pgOrder
}

func (s *PgOrderSuit) TearDownTest() {
	s.cleanup()
}

func (s *PgOrderSuit) TearDownSubTest() {
	s.cleanup()
}

func (s *PgOrderSuit) cleanup() {
	var (
		ctx = context.Background()
		trx = s.transactor.ExtContext(ctx)
	)

	sqlx.MustExecContext(ctx, trx, "DELETE FROM order_status_timeline")
	sqlx.MustExecContext(ctx, trx, "DELETE FROM order_products")
	sqlx.MustExecContext(ctx, trx, "DELETE FROM orders")
}

func connectToDatabase(t *testing.T, driverName string) (*sql.DB, error) {
	t.Helper()

	pool := dockertest.NewPoolT(t, "")

	resource := pool.RunT(t, "postgres",
		dockertest.WithTag("16.3-alpine3.20"),
		dockertest.WithEnv([]string{
			"POSTGRES_DB=bot",
			"POSTGRES_USER=bot",
			"POSTGRES_PASSWORD=bot",
			"listen_addresses = '*'",
		}),
		dockertest.WithoutReuse(),
	)

	var dbConn *sql.DB

	if err := pool.Retry(t.Context(), time.Minute, func() error {
		var err error

		dbConn, err = sql.Open(driverName,
			fmt.Sprintf("host=localhost port=%s user=bot password=bot dbname=bot sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return fmt.Errorf("sql open: %w", err)
		}

		if err := dbConn.PingContext(t.Context()); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	t.Cleanup(func() {
		if err := dbConn.Close(); err != nil {
			t.Errorf("close database connection: %v", err)
		}
	})

	return dbConn, nil
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
