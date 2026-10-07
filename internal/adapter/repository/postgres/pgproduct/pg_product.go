package pgproduct

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/productsvc"
)

var (
	_ productsvc.Repository = (*PgProduct)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgProduct struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgProduct {
	return &PgProduct{
		transactor: transactor,
	}
}
