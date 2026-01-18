package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
)

func (p *Postgres) OutboxSelectForDispatchAnswerPayment(
	ctx context.Context,
	limit int,
) ([]outboxprocessor.OutboxAnswerPayment, error) {
	var (
		query = `
			SELECT
				id,
				payment_id,
				ok,
				error_msg
			FROM
				outbox_answer_payment
			WHERE
				is_dispatched = FALSE
			ORDER BY
				id
			LIMIT
				$1
			FOR UPDATE SKIP LOCKED
		`

		outboxMsgs []model.OutboxAnswerPayment
	)

	if err := sqlx.SelectContext(ctx, p.transactor.ExtContext(ctx), &outboxMsgs, query, limit); err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}

	return model.ToOutboxAnswerPayments(outboxMsgs), nil
}
