package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (p *PgOrder) HistoryOrdersByOffset(
	ctx context.Context,
	chatID msginfo.ChatID,
	offset int,
	limit int,
) ([]order.HistoryOrder, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			ROW_NUMBER() OVER (ORDER BY id) AS serial_number,
			status,
			currency_id,
			created_at,
			total_price
		FROM
			orders
		WHERE
			chat_id = :chat_id
		ORDER BY
			id DESC
		LIMIT
			:limit
		OFFSET
			:offset
	`,
		map[string]any{
			"chat_id": chatID.Int64(),
			"offset":  offset,
			"limit":   limit,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	var (
		dbOrders []model.HistoryOrder
		trx      = p.transactor.ExtContext(ctx)
	)

	if err := sqlx.SelectContext(
		ctx,
		trx,
		&dbOrders,
		trx.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	orders, err := model.ToDomHistoryOrders(dbOrders)
	if err != nil {
		return nil, fmt.Errorf("convert to port ordres: %w", err)
	}

	return orders, nil
}
