package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) GetActiveOrderPosition(
	ctx context.Context,
	orderID order.ID,
) (int, error) {
	pos, err := s.repo.GetOrderPositionByStatus(
		ctx,
		orderID,
		order.StatusConfirmed,
		order.StatusInProgress,
	)

	if err != nil {
		return 0, fmt.Errorf("get position from repo: %w", err)
	}

	return pos, nil
}
