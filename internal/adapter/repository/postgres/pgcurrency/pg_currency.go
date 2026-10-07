package pgcurrency

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/currencysvc"
)

var (
	_ currencysvc.Repository = (*PgCurrency)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgCurrency struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgCurrency {
	return &PgCurrency{
		transactor: transactor,
	}
}
