package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
)

func (p *Postgres) OutboxSelectForDispatchInvoice(
	ctx context.Context,
	limit int,
) ([]outboxprocessor.OutboxInvoice, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				msg_text,
				order_id
			FROM
				outbox_order_invoice
			WHERE
				is_dispatched = FALSE
			ORDER BY
				id
			LIMIT
				$1
			FOR UPDATE SKIP LOCKED
		`

		outboxMsgs []model.OutboxInvoice
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &outboxMsgs, query, limit); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	return model.ToOutboxInvoices(outboxMsgs), nil
}
