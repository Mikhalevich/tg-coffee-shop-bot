package handler

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/app/handler/request"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

type UpdateOrderStatusInput struct {
	ID   string `path:"id" maxLength:"30" example:"123" doc:"Order id"`
	Body request.OrderStatus
}

type UpdateOrderStatusOutput struct {
}

func (h *Handler) UpdateOrderStatus(
	ctx context.Context,
	input *UpdateOrderStatusInput,
) (*UpdateOrderStatusOutput, error) {
	orderID, err := order.IDFromString(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid id format")
	}

	status, err := order.StatusFromString(input.Body.Status)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid status")
	}

	if err := h.updateOrderStatusUsecase.Update(ctx, orderID, status); err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			return nil, huma.Error404NotFound("no orders with relevant status")
		}

		return nil, fmt.Errorf("update order status: %w", err)
	}

	return &UpdateOrderStatusOutput{}, nil
}
