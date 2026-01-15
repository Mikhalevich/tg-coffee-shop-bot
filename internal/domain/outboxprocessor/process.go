package outboxprocessor

import (
	"context"
	"errors"
	"fmt"
)

func (o *OutboxProcessor) Process(ctx context.Context, batchSize int) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := o.repository.OutboxSelectForDispatchMessages(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox messages: %w", err)
		}

		var (
			ids  = make([]int, 0, len(msgs))
			errs error
		)

		for _, msg := range msgs {
			if err := o.sender.SendMessage(ctx, msg.Message); err != nil {
				errs = errors.Join(errs, fmt.Errorf("send message: %w", err))

				continue
			}

			ids = append(ids, msg.ID)
		}

		if len(ids) > 0 {
			if err := o.repository.OutboxSetDispatched(ctx, ids, o.timeProvider.Now()); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return errs
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
