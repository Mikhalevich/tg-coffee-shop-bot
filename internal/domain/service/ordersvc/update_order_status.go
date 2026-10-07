package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (s *Service) UpdateOrderStatus(
	ctx context.Context,
	orderID order.ID,
	newStatus order.Status,
) error {
	prevStatuses, err := calculateLegalPreviousStatuses(newStatus)
	if err != nil {
		return fmt.Errorf("calculate prev statuses: %w", err)
	}

	now := s.timeProvider.Now()

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdateOrderStatus(
			ctx,
			orderID,
			now,
			newStatus,
			prevStatuses...,
		); err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		if err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			newStatus,
			now,
		); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
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
