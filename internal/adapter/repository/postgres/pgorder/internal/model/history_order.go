package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type HistoryOrder struct {
	ID           int       `db:"id"`
	SerialNumber int       `db:"serial_number"`
	Status       string    `db:"status"`
	CurrencyID   int       `db:"currency_id"`
	CreatedAt    time.Time `db:"created_at"`
	TotalPrice   int       `db:"total_price"`
}

func ToDomHistoryOrders(orders []HistoryOrder) ([]order.HistoryOrder, error) {
	shortOrders := make([]order.HistoryOrder, 0, len(orders))

	for _, v := range orders {
		domShortOrder, err := ToDomHistoryOrder(v)
		if err != nil {
			return nil, fmt.Errorf("convert to short order: %w", err)
		}

		shortOrders = append(shortOrders, domShortOrder)
	}

	return shortOrders, nil
}

func ToDomHistoryOrder(dbOrder HistoryOrder) (order.HistoryOrder, error) {
	status, err := order.StatusFromString(dbOrder.Status)
	if err != nil {
		return order.HistoryOrder{}, fmt.Errorf("status from string: %w", err)
	}

	return order.HistoryOrder{
		ID:           order.IDFromInt(dbOrder.ID),
		SerialNumber: dbOrder.SerialNumber,
		Status:       status,
		CurrencyID:   currency.IDFromInt(dbOrder.CurrencyID),
		CreatedAt:    dbOrder.CreatedAt,
		TotalPrice:   dbOrder.TotalPrice,
	}, nil
}
