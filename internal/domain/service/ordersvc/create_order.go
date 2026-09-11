package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) CreateOrder(
	ctx context.Context,
	info order.CreateOrderInfo,
) (order.Order, error) {
	var (
		newOrder = order.Order{
			ChatID:           info.ChatID,
			Status:           info.Status,
			VerificationCode: info.VerificationCode,
			CurrencyID:       info.CurrencyID,
			TotalPrice:       info.TotalPrice,
			Timeline: []order.StatusTime{
				{
					Status: info.Status,
					Time:   info.StatusOperationTime,
				},
			},
			Products:  info.Products,
			CreatedAt: info.StatusOperationTime,
			UpdatedAt: info.StatusOperationTime,
		}
	)

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		orderID, err := s.repo.InsertOrder(
			ctx,
			newOrder,
		)

		if err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		if err := s.repo.InsertProductsOrder(
			ctx,
			orderID,
			newOrder.Products,
		); err != nil {
			return fmt.Errorf("insert products order: %w", err)
		}

		if err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			newOrder.Status,
			newOrder.CreatedAt,
		); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		newOrder.ID = orderID

		return nil
	}); err != nil {
		return order.Order{}, fmt.Errorf("transaction: %w", err)
	}

	return newOrder, nil
}
