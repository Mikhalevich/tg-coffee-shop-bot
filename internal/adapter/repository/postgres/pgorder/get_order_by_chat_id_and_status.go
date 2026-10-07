package pgorder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) GetOrderByChatIDAndStatus(
	ctx context.Context,
	id msginfo.ChatID,
	statuses ...order.Status,
) (order.Order, error) {
	trx := p.transactor.ExtContext(ctx)
	dbOrder, err := selectOrderByChatIDAndStatus(ctx, trx, id, statuses...)
	if err != nil {
		return order.Order{}, fmt.Errorf("select order by chat id and status: %w", err)
	}

	orderProducts, err := selectOrderProducts(ctx, trx, dbOrder.ID)
	if err != nil {
		return order.Order{}, fmt.Errorf("select order products: %w", err)
	}

	orderTimeline, err := selectOrderTimeline(ctx, trx, dbOrder.ID)
	if err != nil {
		return order.Order{}, fmt.Errorf("select order timeline: %w", err)
	}

	domOrder, err := dbOrder.ToDom(orderProducts, orderTimeline)
	if err != nil {
		return order.Order{}, fmt.Errorf("convert to port order: %w", err)
	}

	return domOrder, nil
}

func selectOrderByChatIDAndStatus(
	ctx context.Context,
	ext sqlx.ExtContext,
	chatID msginfo.ChatID,
	statuses ...order.Status,
) (model.Order, error) {
	query, args, err := sqlx.In(`
		SELECT
			id,
			chat_id,
			status,
			verification_code,
			currency_id,
			daily_position,
			total_price,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			chat_id = ? AND
			status IN(?)
	`, chatID, statuses)

	if err != nil {
		return model.Order{}, fmt.Errorf("sqlx in: %w", err)
	}

	var dbOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &dbOrder, ext.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, perror.NotFound("order not found")
		}

		return model.Order{}, fmt.Errorf("get context: %w", err)
	}

	return dbOrder, nil
}
