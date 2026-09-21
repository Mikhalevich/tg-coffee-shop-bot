package pgorder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
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
	query, args, err := sqlx.Named(`
		UPDATE orders SET
			status = :status,
			updated_at = :updated_at
		WHERE
			id = :id AND
			chat_id = :chat_id AND
			status IN (?)
		RETURNING *
		`,
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

	var (
		trx     = p.transactor.ExtContext(ctx)
		dbOrder model.Order
	)

	if err := sqlx.GetContext(
		ctx,
		trx,
		&dbOrder,
		trx.Rebind(query),
		args...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return perror.NoRowsUpdated()
		}

		return fmt.Errorf("get context: %w", err)
	}

	return nil
}
