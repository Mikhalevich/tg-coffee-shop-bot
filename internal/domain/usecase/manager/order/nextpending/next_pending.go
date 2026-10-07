package nextpending

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Transactor interface {
	Transaction(
		ctx context.Context,
		trxFn func(ctx context.Context) error,
	) error
}

type OrderService interface {
	GetNextPendingOrder(ctx context.Context) (order.Order, error)
}

type NotificationService interface {
	SendOrderStatusChanged(
		ctx context.Context,
		chatID msginfo.ChatID,
		orderStatus order.Status,
	) error
}

type NextPending struct {
	transactor          Transactor
	orderService        OrderService
	notificationService NotificationService
}

func New(
	transactor Transactor,
	orderService OrderService,
	notificationService NotificationService,
) *NextPending {
	return &NextPending{
		transactor:          transactor,
		orderService:        orderService,
		notificationService: notificationService,
	}
}

func (n *NextPending) Next(
	ctx context.Context,
) (order.Order, error) {
	var (
		nextOrder order.Order
		err       error
	)

	if err := n.transactor.Transaction(ctx, func(ctx context.Context) error {
		nextOrder, err = n.orderService.GetNextPendingOrder(ctx)
		if err != nil {
			return fmt.Errorf("get next pending order: %w", err)
		}

		if err := n.notificationService.SendOrderStatusChanged(
			ctx,
			nextOrder.ChatID,
			nextOrder.Status,
		); err != nil {
			return fmt.Errorf("order status changed msg: %w", err)
		}

		return nil
	}); err != nil {
		return order.Order{}, fmt.Errorf("transaction: %w", err)
	}

	return nextOrder, nil
}
