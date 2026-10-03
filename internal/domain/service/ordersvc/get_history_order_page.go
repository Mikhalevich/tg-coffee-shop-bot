package ordersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

const (
	lastPageMarker = -1
)

func (s *Service) GetHistoryOrderPage(
	ctx context.Context,
	chatID msginfo.ChatID,
	page int,
	pageSize int,
) ([]order.HistoryOrder, order.Page, error) {
	ordersCount, err := s.repo.HistoryOrdersCount(ctx, chatID)
	if err != nil {
		return nil, order.Page{}, fmt.Errorf("history orders count: %w", err)
	}

	if ordersCount == 0 {
		return nil, order.Page{}, perror.NotFound("orders not found")
	}

	pagesCount := calculatePageCount(ordersCount, pageSize)

	if page == lastPageMarker {
		page = pagesCount
	}

	if page > pagesCount {
		return nil, order.Page{}, fmt.Errorf("invalid page number: %d, pages total: %d", page, pagesCount)
	}

	orders, err := s.repo.HistoryOrdersByOffset(
		ctx,
		chatID,
		calculatePageOffset(page, pageSize),
		pageSize,
	)

	if err != nil {
		return nil, order.Page{}, fmt.Errorf("history orders by offset: %w", err)
	}

	return orders, order.Page{
		Number: page,
		Total:  pagesCount,
	}, nil
}

func calculatePageOffset(pageNumber, pageSize int) int {
	if pageNumber <= 1 {
		return 0
	}

	return (pageNumber - 1) * pageSize
}

func calculatePageCount(count, pageSize int) int {
	var (
		fullPages    = count / pageSize
		lastPageSize = count % pageSize
	)

	if lastPageSize > 0 {
		return fullPages + 1
	}

	return fullPages
}
