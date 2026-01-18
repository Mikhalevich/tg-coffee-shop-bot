package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (p *Postgres) UpdateOrderByChatAndID(
	ctx context.Context,
	orderID order.ID,
	chatID msginfo.ChatID,
	data order.UpdateOrderData,
	prevStatuses ...order.Status,
) (*order.Order, error) {
	var (
		dbOrder       *model.Order
		orderProducts []model.OrderProduct
		orderTimeline []model.OrderTimeline
		err           error
	)

	if err := p.transactor.Transaction(ctx, func(ctx context.Context) error {
		trx := p.transactor.ExtContext(ctx)
		dbOrder, err = updateOrderDataByChatAndID(ctx, trx, orderID, chatID, data, prevStatuses...)
		if err != nil {
			return fmt.Errorf("update order data: %w", err)
		}

		if err := insertOrderTimeline(ctx, trx, model.OrderTimeline{
			ID:        orderID.Int(),
			Status:    data.Status.String(),
			UpdatedAt: data.StatusOperationTime,
		}); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		orderProducts, err = selectOrderProducts(ctx, trx, dbOrder.ID)
		if err != nil {
			return fmt.Errorf("select order products: %w", err)
		}

		orderTimeline, err = selectOrderTimeline(ctx, trx, orderID.Int())
		if err != nil {
			return fmt.Errorf("select order timeline: %w", err)
		}

		return nil
	},
	); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	portOrder, err := model.ToPortOrder(dbOrder, orderProducts, orderTimeline)
	if err != nil {
		return nil, fmt.Errorf("convert to port order: %w", err)
	}

	return portOrder, nil
}

func updateOrderDataByChatAndID(
	ctx context.Context,
	ext sqlx.ExtContext,
	orderID order.ID,
	chatID msginfo.ChatID,
	data order.UpdateOrderData,
	prevStatuses ...order.Status,
) (*model.Order, error) {
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
		RETURNING *
		`,
		map[string]any{
			"status":            data.Status,
			"verification_code": model.NullString(data.VerificationCode),
			//nolint:gosec
			"daily_position": model.NullIntPositive(int32(data.DailyPosition)),
			"updated_at":     data.StatusOperationTime,
			"id":             orderID.Int(),
			"chat_id":        chatID.Int64(),
		})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	args = append(args, prevStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, fmt.Errorf("in statement %w", err)
	}

	var dbOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &dbOrder, ext.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotUpdated
		}

		return nil, fmt.Errorf("get context: %w", err)
	}

	return &dbOrder, nil
}
