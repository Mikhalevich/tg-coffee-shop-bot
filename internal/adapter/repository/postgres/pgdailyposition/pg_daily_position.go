package pgdailyposition

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/payment"
)

var (
	_ payment.PositionService = (*PgDailyPosition)(nil)
)

type Transactor interface {
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type PgDailyPosition struct {
	transactor Transactor
}

func New(
	transactor Transactor,
) *PgDailyPosition {
	return &PgDailyPosition{
		transactor: transactor,
	}
}
