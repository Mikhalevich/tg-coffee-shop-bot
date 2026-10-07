package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/null"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) UpdateOrderByChatAndID(
	ctx context.Context,
	orderID order.ID,
	chatID msginfo.ChatID,
	data order.UpdateOrderData,
	prevStatuses ...order.Status,
) error {
	query, args, err := sqlx.Named(`
		UPDATE orders SET
			status = :status,
			verification_code = :verification_code,
			daily_position = :daily_position,
			updated_at = :updated_at
		WHERE
			id = :id AND
			chat_id = :chat_id AND
			status IN (?)
		`,
		map[string]any{
			"status":            data.Status,
			"verification_code": null.String(data.VerificationCode),
			//nolint:gosec
			"daily_position": null.IntPositive(int32(data.DailyPosition)),
			"updated_at":     data.StatusOperationTime,
			"id":             orderID.Int(),
			"chat_id":        chatID.Int64(),
		})

	if err != nil {
		return fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return fmt.Errorf("in statement %w", err)
	}

	trx := p.transactor.ExtContext(ctx)

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
