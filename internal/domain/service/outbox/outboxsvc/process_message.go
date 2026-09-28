package outboxsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (s *Service) ProcessMessage(ctx context.Context, batchSize int) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := s.repo.SelectForDispatchMessages(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox messages: %w", err)
		}

		var (
			ids  = make([]int, 0, len(msgs))
			errs error
		)

		for _, msg := range msgs {
			logger.FromContext(ctx).
				WithFields(
					logger.Fields{
						"chat_id":  msg.ChatID,
						"text":     msg.Text,
						"msg_type": msg.Type,
					},
				).
				Debug("send message")

			if err := s.sender.SendMessage(ctx, msg.Message); err != nil {
				errs = errors.Join(errs, fmt.Errorf("send message: %w", err))

				continue
			}

			ids = append(ids, msg.ID)
		}

		if len(ids) > 0 {
			if err := s.repo.SetDispatched(ctx, ids, s.timeProvider.Now()); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return errs
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
