package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) GetNextPendingOrder(
	ctx context.Context,
) (order.Order, error) {
	var (
		ord order.Order
	)
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		now := s.timeProvider.Now()

		orderID, err := s.repo.UpdateOrderStatusForMinID(
			ctx,
			now,
			order.StatusInProgress,
			order.StatusConfirmed,
		)

		if err != nil {
			return fmt.Errorf("update order status for min id: %w", err)
		}

		if err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			order.StatusInProgress,
			now,
		); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		ord, err = s.repo.GetOrderByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("get order by id: %w", err)
		}

		return nil
	}); err != nil {
		return order.Order{}, fmt.Errorf("transaction: %w", err)
	}

	return ord, nil
}
