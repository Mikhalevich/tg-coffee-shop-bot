package queuesize

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type OrderService interface {
	GetOrdersCount(
		ctx context.Context,
	) (int, error)
}

type NotificationService interface {
	ShowOrderCount(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		count int,
	) error
}

type QueueSize struct {
	orderService        OrderService
	notificationService NotificationService
}

func New(
	orderService OrderService,
	notificationService NotificationService,
) *QueueSize {
	return &QueueSize{
		orderService:        orderService,
		notificationService: notificationService,
	}
}

func (q *QueueSize) Size(
	ctx context.Context,
	info msginfo.Info,
) error {
	count, err := q.orderService.GetOrdersCount(ctx)
	if err != nil {
		return fmt.Errorf("get order count: %w", err)
	}

	if err := q.notificationService.ShowOrderCount(
		ctx,
		info.ChatID,
		info.MessageID,
		count,
	); err != nil {
		return fmt.Errorf("show order count msg: %w", err)
	}

	return nil
}
