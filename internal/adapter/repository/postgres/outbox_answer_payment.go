package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/model"
)

func (p *Postgres) AnswerOrderPayment(
	ctx context.Context,
	paymentID string,
	paymentOK bool,
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

		answerPayment = model.OutboxAnswerPayment{
			PaymentID: paymentID,
			OK:        paymentOK,
			ErrorMsg:  errorMsg,
		}
	)

	res, err := sqlx.NamedExecContext(ctx, p.Transactor().ExtContext(ctx), query, answerPayment)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
