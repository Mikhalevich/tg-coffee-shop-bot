package pgproduct_test

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
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgproduct"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func productCreatedAt() time.Time {
	return time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
}

func productUpdatedAt() time.Time {
	return time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC)
}

type ProductSuit struct {
	*suite.Suite

	transactor pgproduct.Transactor
	pgProduct  *pgproduct.PgProduct
}

func TestProductSuit(t *testing.T) {
	t.Parallel()

	suite.Run(t, &ProductSuit{
		Suite: new(suite.Suite),
	})
}

func (s *ProductSuit) SetupSuite() {
	dbDriver := driver.NewPgx()

	dbConn, err := connectToDatabase(s.T(), dbDriver.Name())
	if err != nil {
		s.FailNow("could not connect to database", err)
	}

	if err := migrationsUp(dbConn, "../../../../../script/db/migrations"); err != nil {
		s.FailNow("could not exec migrations", err)
	}

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, dbDriver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		pgProduct           = pgproduct.New(transactionProvider)
	)

	s.transactor = transactionProvider
	s.pgProduct = pgProduct
}

func (s *ProductSuit) TearDownTest() {
	s.cleanup()
}

func (s *ProductSuit) TearDownSubTest() {
	s.cleanup()
}

func (s *ProductSuit) cleanup() {
	var (
		ctx = context.Background()
	)

	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM product_category")
	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM product_price")
	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM category")
	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM product")
	sqlx.MustExecContext(ctx, s.transactor.ExtContext(ctx), "DELETE FROM currency")
}

func (s *ProductSuit) insertCategory(title string, isEnabled bool) int {
	var (
		ctx        = s.T().Context()
		categoryID int
	)

	err := sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &categoryID, `
		INSERT INTO category (
			title,
			is_enabled
		) VALUES (
			$1,
			$2
		) RETURNING id`,
		title,
		isEnabled,
	)
	s.Require().NoError(err)

	return categoryID
}

func (s *ProductSuit) insertProduct(title string, isEnabled bool) int {
	var (
		ctx       = s.T().Context()
		productID int
	)

	err := sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &productID, `
		INSERT INTO product (
			title,
			is_enabled,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id`,
		title,
		isEnabled,
		productCreatedAt(),
		productUpdatedAt(),
	)
	s.Require().NoError(err)

	return productID
}

func (s *ProductSuit) linkProductCategory(productID int, categoryID int) {
	ctx := s.T().Context()

	_, err := s.transactor.ExtContext(ctx).ExecContext(ctx, `
		INSERT INTO product_category (
			product_id,
			category_id
		) VALUES (
			$1,
			$2
		)`,
		productID,
		categoryID,
	)
	s.Require().NoError(err)
}

func (s *ProductSuit) insertCurrency(code string) int {
	var (
		ctx        = s.T().Context()
		currencyID int
	)

	err := sqlx.GetContext(ctx, s.transactor.ExtContext(ctx), &currencyID, `
		INSERT INTO currency (
			code,
			exp,
			decimal_sep,
			min_amount,
			max_amount,
			is_enabled
		) VALUES (
			$1,
			2,
			'.',
			100,
			100000,
			TRUE
		) RETURNING id`,
		code,
	)
	s.Require().NoError(err)

	return currencyID
}

func (s *ProductSuit) insertPrice(productID int, currencyID int, price int) {
	ctx := s.T().Context()

	_, err := s.transactor.ExtContext(ctx).ExecContext(ctx, `
		INSERT INTO product_price (
			product_id,
			currency_id,
			price
		) VALUES (
			$1,
			$2,
			$3
		)`,
		productID,
		currencyID,
		price,
	)
	s.Require().NoError(err)
}

func makeProduct(productID int, title string, currencyID int, price int, isEnabled bool) product.Product {
	return product.Product{
		ID:         product.ProductIDFromInt(productID),
		Title:      title,
		CurrencyID: currency.IDFromInt(currencyID),
		Price:      price,
		IsEnabled:  isEnabled,
		CreatedAt:  productCreatedAt(),
		UpdatedAt:  productUpdatedAt(),
	}
}

// normalizeTime converts timestamps to UTC so products read from postgres can be compared with expected ones.
func normalizeTime(p product.Product) product.Product {
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()

	return p
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
