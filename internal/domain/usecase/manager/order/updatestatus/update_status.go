package updatestatus

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
	GetOrderByID(
		ctx context.Context,
		id order.ID,
	) (order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		orderID order.ID,
		newStatus order.Status,
	) error
}

type NotificationService interface {
	SendOrderStatusChanged(
		ctx context.Context,
		chatID msginfo.ChatID,
		newStatus order.Status,
	) error
}

type UpdateStatus struct {
	transactor          Transactor
	orderService        OrderService
	notificationService NotificationService
}

func New(
	transactor Transactor,
	orderService OrderService,
	notificationService NotificationService,
) *UpdateStatus {
	return &UpdateStatus{
		transactor:          transactor,
		orderService:        orderService,
		notificationService: notificationService,
	}
}

func (u *UpdateStatus) Update(
	ctx context.Context,
	orderID order.ID,
	status order.Status,
) error {
	if err := u.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := u.orderService.UpdateOrderStatus(
			ctx,
			orderID,
			status,
		); err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		ord, err := u.orderService.GetOrderByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("get order: %w", err)
		}

		if err := u.notificationService.SendOrderStatusChanged(
			ctx,
			ord.ChatID,
			status,
		); err != nil {
			return fmt.Errorf("order status changed msg: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
