package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) GetActiveOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
) (order.Order, error) {
	activeOrder, err := s.repo.GetOrderByChatIDAndStatus(
		ctx,
		chatID,
		order.StatusWaitingPayment,
		order.StatusPaymentInProgress,
		order.StatusConfirmed,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		return order.Order{}, fmt.Errorf("get active order from repo: %w", err)
	}

	return activeOrder, nil
}
