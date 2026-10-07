package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (p *PgOrder) GetOrdersCountByStatus(
	ctx context.Context,
	statuses ...order.Status,
) (int, error) {
	query, args, err := sqlx.In(`
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE
			status IN(?)
	`, statuses)

	if err != nil {
		return 0, fmt.Errorf("sqlx in statement: %w", err)
	}

	var (
		trx   = p.transactor.ExtContext(ctx)
		count int
	)

	if err := sqlx.GetContext(
		ctx,
		trx,
		&count,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return count, nil
}
