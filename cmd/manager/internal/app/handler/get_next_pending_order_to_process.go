package handler

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/app/handler/response"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

type GetNextOrderInput struct {
}

type GetNextOrderOutput struct {
	Body response.Order
}

func (h *Handler) GetNextOrder(ctx context.Context, input *GetNextOrderInput) (*GetNextOrderOutput, error) {
	inProgressOrder, err := h.getNextOrderUsecase.Next(ctx)
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			return nil, huma.Error404NotFound("no active orders")
		}

		return nil, fmt.Errorf("GetNextPendingOrderToProcess: %w", err)
	}

	return &GetNextOrderOutput{
		Body: *response.ToOrder(inProgressOrder),
	}, nil
}
