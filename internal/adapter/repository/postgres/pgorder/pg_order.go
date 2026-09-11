package pgorder

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/ordersvc"
)

var (
	_ ordersvc.Repository = (*PgOrder)(nil)
)

type Driver interface {
	IsConstraintError(err error, constraint string) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgOrder struct {
	driver     Driver
	transactor Transactor
}

func New(
	driver Driver,
	transactor Transactor,
) *PgOrder {
	return &PgOrder{
		driver:     driver,
		transactor: transactor,
	}
}
