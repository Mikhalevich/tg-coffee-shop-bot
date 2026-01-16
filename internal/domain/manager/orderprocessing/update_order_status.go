package orderprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (o *OrderProcessing) UpdateOrderStatus(ctx context.Context, orderID order.ID, status order.Status) error {
	previousStatuses, err := calculateLegalPreviousStatuses(status)
	if err != nil {
		return fmt.Errorf("calculate legal previous statuses: %w", err)
	}

	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		updatedOrder, err := o.repository.UpdateOrderStatus(
			ctx,
			orderID,
			o.timeProvider.Now(),
			status,
			previousStatuses...,
		)

		if err != nil {
			if o.repository.IsNotUpdatedError(err) {
				return perror.NotFound("order with relevant status not found")
			}

			return fmt.Errorf("update order status: %w", err)
		}

		if err := o.customerSender.SendMessage(
			ctx,
			messageprocessor.Message{
				ChatID: updatedOrder.ChatID,
				Text:   o.makeChangedOrderStatusMarkdownMsg(status),
				Type:   messageprocessor.MessageTypeMarkdown,
			},
		); err != nil {
			return fmt.Errorf("send outbox message: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func calculateLegalPreviousStatuses(s order.Status) ([]order.Status, error) {
	switch s {
	case order.StatusWaitingPayment:
		return nil, perror.InvalidParam("invalid order transition")
	case order.StatusPaymentInProgress:
		return []order.Status{order.StatusWaitingPayment}, nil
	case order.StatusConfirmed:
		return []order.Status{order.StatusPaymentInProgress}, nil
	case order.StatusInProgress:
		return []order.Status{order.StatusConfirmed}, nil
	case order.StatusReady:
		return []order.Status{order.StatusInProgress}, nil
	case order.StatusCompleted:
		return []order.Status{order.StatusConfirmed, order.StatusInProgress, order.StatusReady}, nil
	case order.StatusCanceled:
		return []order.Status{order.StatusConfirmed, order.StatusInProgress, order.StatusReady}, nil
	case order.StatusRejected:
		return []order.Status{order.StatusConfirmed, order.StatusInProgress, order.StatusReady}, nil
	}

	return nil, perror.InvalidParam("invalid order status")
}
