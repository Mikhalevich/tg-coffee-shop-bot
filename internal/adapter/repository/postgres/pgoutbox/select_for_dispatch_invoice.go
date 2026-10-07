package pgoutbox

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
)

func (p *PgOutbox) SelectForDispatchInvoice(
	ctx context.Context,
	limit int,
) ([]outboxmsg.Invoice, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				title,
				description,
				currency_code,
				order_id,
				labels,
				buttons
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

		outboxMsgs []model.Invoice
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&outboxMsgs,
		query,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	outboxInvoices, err := model.ToOutboxInvoices(outboxMsgs)
	if err != nil {
		return nil, fmt.Errorf("convert to outbox invoices: %w", err)
	}

	return outboxInvoices, nil
}
