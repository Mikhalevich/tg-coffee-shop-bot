package pgproduct

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/productsvc"
)

var (
	_ productsvc.Repository = (*PgProduct)(nil)
)

// Transactor provides the database executor for the current context: the active
// transaction if one is stored in ctx, otherwise the plain connection.
type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

// PgProduct is a PostgreSQL repository for products, categories and product prices.
type PgProduct struct {
	transactor Transactor
}

// New creates a PgProduct that runs its queries through transactor.
func New(
	transactor Transactor,
) *PgProduct {
	return &PgProduct{
		transactor: transactor,
	}
}
