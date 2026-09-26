package pgoutbox

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox/internal/model"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (p *PgOutbox) AnswerPayment(
	ctx context.Context,
	paymentID string,
	success bool,
	errorMsg string,
) error {
	var (
		query = `
			INSERT INTO outbox_answer_payment(
				payment_id,
				ok,
				error_msg
			) VALUES (
				:payment_id,
				:ok,
				:error_msg
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.AnswerPayment{
			PaymentID: paymentID,
			OK:        success,
			ErrorMsg:  errorMsg,
		},
	)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
