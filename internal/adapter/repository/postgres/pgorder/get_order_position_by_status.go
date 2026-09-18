package pgorder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) GetOrderPositionByStatus(
	ctx context.Context,
	orderID order.ID,
	statuses ...order.Status,
) (int, error) {
	query, args, err := sqlx.In(`
		WITH order_queue AS (
			SELECT
				id,
				ROW_NUMBER() OVER (ORDER BY updated_at::date, daily_position) AS position
			FROM
				orders
			WHERE
				status IN(?)
			ORDER BY
				updated_at::date,
				daily_position
		)
		SELECT
			position
		FROM
			order_queue
		WHERE
			id = ?
	`, statuses, orderID.Int())

	if err != nil {
		return 0, fmt.Errorf("sqlx in statement: %w", err)
	}

	var (
		trx = p.transactor.ExtContext(ctx)
		pos int
	)

	if err := sqlx.GetContext(ctx, trx, &pos, trx.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, perror.NotFound("order position not found")
		}

		return 0, fmt.Errorf("get context: %w", err)
	}

	return pos, nil
}
