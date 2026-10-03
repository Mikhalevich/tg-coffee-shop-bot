package pgorder

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (p *PgOrder) HistoryOrdersCount(
	ctx context.Context,
	chatID msginfo.ChatID,
) (int, error) {
	query, args, err := sqlx.Named(`
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE
			chat_id = :chat_id
	`,
		map[string]any{
			"chat_id": chatID.Int64(),
		},
	)

	if err != nil {
		return 0, fmt.Errorf("prepare query: %w", err)
	}

	var (
		count int
		trx   = p.transactor.ExtContext(ctx)
	)

	if err := sqlx.GetContext(
		ctx,
		trx,
		&count,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return count, nil
}
