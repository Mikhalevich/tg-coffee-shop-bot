package pgorder

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) UpdateOrderStatusByChatAndID(
	ctx context.Context,
	orderID order.ID,
	chatID msginfo.ChatID,
	operationTime time.Time,
	newStatus order.Status,
	prevStatuses ...order.Status,
) error {
	var (
		query = `
			UPDATE orders SET
				status = :status,
				updated_at = :updated_at
			WHERE
				id = :id AND
				chat_id = :chat_id AND
				status IN (?)
		`
		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(
		query,
		map[string]any{
			"status":     newStatus,
			"updated_at": operationTime,
			"id":         orderID.Int(),
			"chat_id":    chatID.Int64(),
		})

	if err != nil {
		return fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return fmt.Errorf("in statement %w", err)
	}

	res, err := trx.ExecContext(
		ctx,
		trx.Rebind(query),
		args...,
	)

	if err != nil {
		return fmt.Errorf("get context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
