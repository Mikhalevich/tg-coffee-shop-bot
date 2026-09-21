package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (s *Service) Cancel(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
) error {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if !ord.IsSameChat(chatID) {
		return perror.InvalidParam("order is not beloing to chat_id")
	}

	if !ord.CanCancel() {
		return perror.UnableToCancel("order unable to cancel")
	}

	now := s.timeProvider.Now()

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdateOrderStatusByChatAndID(
			ctx,
			orderID,
			chatID,
			now,
			order.StatusCanceled,
			order.StatusWaitingPayment,
			order.StatusConfirmed,
		); err != nil {
			return fmt.Errorf("udpate order status: %w", err)
		}

		if err := s.repo.InsertOrderTimeline(
			ctx,
			orderID,
			order.StatusCanceled,
			now,
		); err != nil {
			return fmt.Errorf("insert order timeline: %w", err)
		}

		return nil
	}); err != nil {
		if perror.IsType(err, perror.TypeNoRowsUpdated) {
			return perror.UnableToCancel("order unable to cancel")
		}

		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
