package pgorder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) GetOrderByID(
	ctx context.Context,
	id order.ID,
) (order.Order, error) {
	trx := p.transactor.ExtContext(ctx)
	dbOrder, err := selectOrderByID(ctx, trx, id)
	if err != nil {
		return order.Order{}, fmt.Errorf("select order by id: %w", err)
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

func selectOrderByID(
	ctx context.Context,
	ext sqlx.ExtContext,
	orderID order.ID,
) (model.Order, error) {
	var modelOrder model.Order
	if err := sqlx.GetContext(ctx, ext, &modelOrder,
		`SELECT
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
			id = $1
	`, orderID.Int(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, perror.NotFound("order not found")
		}

		return model.Order{}, fmt.Errorf("get context: %w", err)
	}

	return modelOrder, nil
}

func selectOrderProducts(
	ctx context.Context,
	ext sqlx.ExtContext,
	orderID int,
) ([]model.OrderProduct, error) {
	query, args, err := sqlx.Named(`
		SELECT
			order_id,
			product_id,
			count,
			price
		FROM
			order_products
		WHERE
			order_id = :order_id
	`, map[string]any{
		"order_id": orderID,
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var orderProducts []model.OrderProduct
	if err := sqlx.SelectContext(
		ctx,
		ext,
		&orderProducts,
		ext.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return orderProducts, nil
}

func selectOrderTimeline(
	ctx context.Context,
	ext sqlx.ExtContext,
	orderID int,
) ([]model.OrderTimeline, error) {
	query, args, err := sqlx.Named(`
		SELECT
			order_id,
			status,
			updated_at
		FROM
			order_status_timeline
		WHERE
			order_id = :order_id
	`, map[string]any{
		"order_id": orderID,
	})

	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var orderTimeline []model.OrderTimeline
	if err := sqlx.SelectContext(
		ctx,
		ext,
		&orderTimeline,
		ext.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return orderTimeline, nil
}
