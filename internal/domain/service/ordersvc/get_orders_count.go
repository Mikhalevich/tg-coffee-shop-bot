package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) GetOrdersCount(
	ctx context.Context,
) (int, error) {
	size, err := s.repo.GetOrdersCountByStatus(
		ctx,
		order.StatusConfirmed,
		order.StatusInProgress,
	)

	if err != nil {
		return 0, fmt.Errorf("get orders count by status: %w", err)
	}

	return size, nil
}
