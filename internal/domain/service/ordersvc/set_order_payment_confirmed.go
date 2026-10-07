package ordersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) SetOrderPaymentConfirmed(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
	verificationCode string,
	dailyPosition int,
	operationTime time.Time,
) (int, error) {
	var pos int
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdateOrderByChatAndID(
			ctx,
			orderID,
			chatID,
			order.UpdateOrderData{
				Status:              order.StatusConfirmed,
				StatusOperationTime: operationTime,
				VerificationCode:    verificationCode,
				DailyPosition:       dailyPosition,
			},
			order.StatusPaymentInProgress,
		); err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			order.StatusPaymentInProgress,
			operationTime,
		)

		if err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		pos, err = s.repo.GetOrderPositionByStatus(
			ctx,
			orderID,
			order.StatusConfirmed,
		)

		if err != nil {
			return fmt.Errorf("get order position: %w", err)
		}

		return nil
	}); err != nil {
		return 0, fmt.Errorf("transaction: %w", err)
	}

	return pos, nil
}
