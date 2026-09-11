package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOrder) InsertOrder(
	ctx context.Context,
	ord order.Order,
) (order.ID, error) {
	var (
		query = `
			INSERT INTO orders(
				chat_id,
				status,
				verification_code,
				currency_id,
				total_price,
				created_at,
				updated_at
			) VALUES (
				:chat_id,
				:status,
				:verification_code,
				:currency_id,
				:total_price,
				:created_at,
				:updated_at
			)
			RETURNING id
		`

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, model.ToDBOrder(ord))
	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var orderID int
	if err := sqlx.GetContext(
		ctx,
		trx,
		&orderID,
		trx.Rebind(query),
		args...,
	); err != nil {
		if p.driver.IsConstraintError(err, "orders_only_one_active_order_unique_idx") {
			return 0, perror.AlreadyExists("order already exists")
		}

		return 0, fmt.Errorf("insert order: %w", err)
	}

	return order.IDFromInt(orderID), nil
}
