package ordercancel

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

type OrderService interface {
	Cancel(
		ctx context.Context,
		chatID msginfo.ChatID,
		orderID order.ID,
	) error
}

type NotificationService interface {
	InvalidOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	UnableToCancelOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	OrderCanceled(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
}

type OrderCancel struct {
	orderService        OrderService
	notificationService NotificationService
}

func New(
	orderService OrderService,
	notificationService NotificationService,
) *OrderCancel {
	return &OrderCancel{
		orderService:        orderService,
		notificationService: notificationService,
	}
}

func (o *OrderCancel) Cancel(
	ctx context.Context,
	info msginfo.Info,
	orderID order.ID,
) error {
	if err := o.orderService.Cancel(
		ctx,
		info.ChatID,
		orderID,
	); err != nil {
		switch {
		case perror.IsType(err, perror.TypeNotFound):
			if err := o.notificationService.InvalidOrder(ctx, info.ChatID); err != nil {
				return fmt.Errorf("invalid order msg: %w", err)
			}

		case perror.IsType(err, perror.TypeUnableToCancel):
			if err := o.notificationService.UnableToCancelOrder(ctx, info.ChatID); err != nil {
				return fmt.Errorf("unable to cancel order msg: %w", err)
			}
		}

		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}
