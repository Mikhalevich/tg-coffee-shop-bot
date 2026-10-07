package history

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (h *History) First(
	ctx context.Context,
	info msginfo.Info,
) error {
	if err := h.loadPage(ctx, info, 1); err != nil {
		return fmt.Errorf("load fist page: %w", err)
	}

	return nil
}

func (h *History) loadPage(
	ctx context.Context,
	info msginfo.Info,
	pageNumber int,
) error {
	orders, pageInfo, err := h.orderService.GetHistoryOrderPage(ctx, info.ChatID, pageNumber, h.pageSize)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("get history orders page: %w", err)
		}

		if err := h.notificationService.NoOrdersFound(ctx, info.ChatID); err != nil {
			return fmt.Errorf("no orders notification: %w", err)
		}

		return nil
	}

	curr, err := h.currencyService.GetCurrencyByID(ctx, orders[0].CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := h.notificationService.ShowHistoryPage(
		ctx,
		info.ChatID,
		info.MessageID,
		orders,
		pageInfo,
		curr,
	); err != nil {
		return fmt.Errorf("show history page: %w", err)
	}

	return nil
}
