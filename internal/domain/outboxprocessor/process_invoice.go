package outboxprocessor

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (o *OutboxProcessor) ProcessInvoice(ctx context.Context, batchSize int) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		msgs, err := o.repository.OutboxSelectForDispatchInvoice(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("select outbox invoices: %w", err)
		}

		var (
			ids  = make([]int, 0, len(msgs))
			errs error
		)

		for _, msg := range msgs {
			if err := o.processInvoice(ctx, msg); err != nil {
				errs = errors.Join(errs, fmt.Errorf("process invoice: %w", err))

				continue
			}

			ids = append(ids, msg.ID)
		}

		if len(ids) > 0 {
			if err := o.repository.OutboxSetInvoiceDispatched(ctx, ids, o.timeProvider.Now()); err != nil {
				return fmt.Errorf("set dispatched: %w", err)
			}
		}

		return errs
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (o *OutboxProcessor) processInvoice(ctx context.Context, msg OutboxInvoice) error {
	logger.FromContext(ctx).
		WithFields(
			logger.Fields{
				"chat_id":  msg.ChatID.Int64(),
				"text":     msg.Text,
				"order_id": msg.OrderID.Int(),
			},
		).
		Debug("send invoice")

	ord, err := o.repository.GetOrderByID(ctx, msg.OrderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	curr, err := o.repository.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	productInfo, err := o.repository.GetProductsByIDs(ctx, ord.ProductIDs(), ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	if err := o.sender.SendInvoice(
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
