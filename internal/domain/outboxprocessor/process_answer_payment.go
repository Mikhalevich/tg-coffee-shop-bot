package outboxprocessor

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (o *OutboxProcessor) ProcessAnswerPayment(ctx context.Context, batchSize int) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		payments, err := o.repository.OutboxSelectForDispatchAnswerPayment(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox answer payment: %w", err)
		}

		var (
			ids  = make([]int, 0, len(payments))
			errs error
		)

		for _, payment := range payments {
			logger.FromContext(ctx).
				WithFields(
					logger.Fields{
						"payment_id": payment.PaymentID,
						"ok":         payment.OK,
						"error_msg":  payment.ErrorMsg,
					},
				).
				Debug("answer order payment")

			if err := o.sender.AnswerOrderPayment(
				ctx,
				payment.PaymentID,
				payment.OK,
				payment.ErrorMsg,
			); err != nil {
				errs = errors.Join(errs, fmt.Errorf("answer payment: %w", err))

				continue
			}

			ids = append(ids, payment.ID)
		}

		if len(ids) > 0 {
			if err := o.repository.OutboxSetAnswerPaymentDispatched(ctx, ids, o.timeProvider.Now()); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return errs
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
