package outboxsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (s *Service) ProcessInvoice(ctx context.Context, batchSize int) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := s.repo.SelectForDispatchInvoice(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox invoices: %w", err)
		}

		var (
			ids  = make([]int, 0, len(msgs))
			errs error
		)

		for _, msg := range msgs {
			if err := s.processInvoice(ctx, msg); err != nil {
				errs = errors.Join(errs, fmt.Errorf("process invoice: %w", err))

				continue
			}

			ids = append(ids, msg.ID)
		}

		if len(ids) > 0 {
			if err := s.repo.SetInvoiceDispatched(
				ctx,
				ids,
				s.timeProvider.Now(),
			); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return errs
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (s *Service) processInvoice(ctx context.Context, msg outboxmsg.Invoice) error {
	logger.FromContext(ctx).
		WithFields(
			logger.Fields{
				"chat_id":  msg.ChatID.Int64(),
				"text":     msg.Text,
				"order_id": msg.OrderID.Int(),
			},
		).
		Debug("send invoice")

	ord, err := s.orderService.GetOrderByID(ctx, msg.OrderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	curr, err := s.currencyService.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	productInfo, err := s.productsService.GetProductsByIDs(
		ctx,
		ord.ProductIDs(),
		ord.CurrencyID,
	)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	if err := s.sender.SendInvoice(
		ctx,
		msg.ChatID,
		msg.Text,
		ord,
		productInfo,
		curr,
	); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}
