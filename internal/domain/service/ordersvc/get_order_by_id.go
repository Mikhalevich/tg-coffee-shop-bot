package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) GetOrderByID(
	ctx context.Context,
	id order.ID,
) (order.Order, error) {
	ord, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return order.Order{}, fmt.Errorf("get order: %w", err)
	}

	return ord, nil
}
