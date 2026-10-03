package pgorder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) UpdateOrderStatusForMinID(
	ctx context.Context,
	operationTime time.Time,
	newStatus, prevStatus order.Status,
) (order.ID, error) {
	query, args, err := sqlx.Named(`
		UPDATE orders SET
			status = :new_status,
			updated_at = :updated_at
		WHERE id = (
				SELECT MIN(id)
				FROM orders
				WHERE status = :previous_status
			)
		RETURNING id
		`, map[string]any{
		"new_status":      newStatus,
		"updated_at":      operationTime,
		"previous_status": prevStatus,
	})

	if err != nil {
		return 0, fmt.Errorf("named: %w", err)
	}

	var (
		orderID int
		trx     = p.transactor.ExtContext(ctx)
	)

	if err := sqlx.GetContext(
		ctx,
		trx, &orderID,
		trx.Rebind(query),
		args...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, perror.NotFound("order not found")
		}

		return 0, fmt.Errorf("get context: %w", err)
	}

	return order.IDFromInt(orderID), nil
}
