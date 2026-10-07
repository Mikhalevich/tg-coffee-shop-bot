package handler

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type GetNextOrderUsecase interface {
	Next(
		ctx context.Context,
	) (order.Order, error)
}

type UpdateOrderStatusUsecase interface {
	Update(
		ctx context.Context,
		orderID order.ID,
		status order.Status,
	) error
}

type Handler struct {
	getNextOrderUsecase      GetNextOrderUsecase
	updateOrderStatusUsecase UpdateOrderStatusUsecase
}

func New(
	getNextOrderUsecase GetNextOrderUsecase,
	updateOrderStatusUsecase UpdateOrderStatusUsecase,
) *Handler {
	return &Handler{
		getNextOrderUsecase:      getNextOrderUsecase,
		updateOrderStatusUsecase: updateOrderStatusUsecase,
	}
}
