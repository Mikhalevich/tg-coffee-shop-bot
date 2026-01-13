package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) OutboxSelectForDispatchMessages(
	ctx context.Context,
	limit int,
) ([]outboxprocessor.OutboxMessage, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				reply_msg_id,
				msg_text,
				payload,
				buttons
			FROM
				outbox_messages
			WHERE
				is_dispatched = FALSE
			ORDER BY
				id
			LIMIT
				$1
			FOR UPDATE SKIP LOCKED
		`

		messages []model.OutboxMessage
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &messages, query, limit); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	return model.ToOutboxMessages(messages), nil
}
