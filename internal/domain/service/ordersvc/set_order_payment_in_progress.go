package ordersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (s *Service) SerOrderPaymentInProgress(
	ctx context.Context,
	orderID order.ID,
	totalAmount int,
	operationTime time.Time,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		ord, err := s.repo.GetOrderByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("get order by id: %w", err)
		}

		if ord.Status != order.StatusWaitingPayment {
			return perror.InvalidParam("invalid order status")
		}

		if ord.TotalPrice != totalAmount {
			return perror.InvalidParam("invalid order total amount")
		}

		if err := s.repo.UpdateOrderStatus(
			ctx,
			orderID,
			operationTime,
			order.StatusPaymentInProgress,
			order.StatusWaitingPayment,
		); err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		if err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			order.StatusPaymentInProgress,
			operationTime,
		); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
