package orderprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (o *OrderProcessing) GetNextPendingOrderToProcess(ctx context.Context) (*order.Order, error) {
	var (
		orderToProcess *order.Order
		err            error
	)
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		orderToProcess, err = o.repository.UpdateOrderStatusForMinID(
			ctx,
			o.timeProvider.Now(),
			order.StatusInProgress,
			order.StatusConfirmed,
		)
		if err != nil {
			if o.repository.IsNotFoundError(err) {
				return perror.NotFound("no pending orders")
			}

			return fmt.Errorf("update next order status: %w", err)
		}

		if err := o.customerSender.SendMessage(
			ctx,
			messageprocessor.Message{
				ChatID: orderToProcess.ChatID,
				Text:   o.makeChangedOrderStatusMarkdownMsg(orderToProcess.Status),
				Type:   messageprocessor.MessageTypeMarkdown,
			},
		); err != nil {
			return fmt.Errorf("send outbox message: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %w", err)
	}

	return orderToProcess, nil
}

func (o *OrderProcessing) makeChangedOrderStatusMarkdownMsg(s order.Status) string {
	return fmt.Sprintf("your order status changed to *%s*", o.escaper.EscapeMarkdown(s.HumanReadable()))
}
