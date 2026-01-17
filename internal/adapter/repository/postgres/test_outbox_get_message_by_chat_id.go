package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
)

func (p *Postgres) TestGetOutboxMessageByChatID(ctx context.Context, chatID int) (model.OutboxMessage, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				reply_msg_id,
				msg_text,
				msg_type,
				payload,
				buttons,
				is_dispatched,
				created_at,
				dispatched_at
			FROM
				outbox_messages
			WHERE
				chat_id = :chat_id
			LIMIT
				1
		`

		msg model.OutboxMessage
	)

	query, args, err := sqlx.Named(query, map[string]any{
		"chat_id": chatID,
	})

	if err != nil {
		return msg, fmt.Errorf("sqlx named: %w", err)
	}

	trx := p.transactor.ExtContext(ctx)

	if err := sqlx.GetContext(ctx, trx, &msg, trx.Rebind(query), args...); err != nil {
		return msg, fmt.Errorf("get context: %w", err)
	}

	return msg, nil
}
