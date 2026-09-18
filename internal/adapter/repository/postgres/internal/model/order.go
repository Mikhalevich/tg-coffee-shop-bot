package model

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type Order struct {
	ID               int            `db:"id"`
	ChatID           int64          `db:"chat_id"`
	Status           string         `db:"status"`
	VerificationCode sql.NullString `db:"verification_code"`
	CurrencyID       int            `db:"currency_id"`
	DailyPosition    sql.NullInt32  `db:"daily_position"`
	TotalPrice       int            `db:"total_price"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

type OrderTimeline struct {
	ID        int       `db:"order_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrderProduct struct {
	OrderID   int `db:"order_id"`
	ProductID int `db:"product_id"`
	Count     int `db:"count"`
	Price     int `db:"price"`
}

type HistoryOrder struct {
	ID           int       `db:"id"`
	SerialNumber int       `db:"serial_number"`
	Status       string    `db:"status"`
	CurrencyID   int       `db:"currency_id"`
	CreatedAt    time.Time `db:"created_at"`
	TotalPrice   int       `db:"total_price"`
}

func PortToOrderProducts(id order.ID, portProducts []order.OrderedProduct) []OrderProduct {
	dbProducts := make([]OrderProduct, 0, len(portProducts))

	for _, v := range portProducts {
		dbProducts = append(dbProducts, OrderProduct{
			OrderID:   id.Int(),
			ProductID: v.ProductID.Int(),
			Count:     v.Count,
			Price:     v.Price,
		})
	}

	return dbProducts
}

func toPortCartProducts(dbProducts []OrderProduct) []order.OrderedProduct {
	portProducts := make([]order.OrderedProduct, 0, len(dbProducts))

	for _, v := range dbProducts {
		portProducts = append(portProducts, order.OrderedProduct{
			ProductID: product.ProductIDFromInt(v.ProductID),
			Count:     v.Count,
			Price:     v.Price,
		})
	}

	return portProducts
}

func ToPortHistoryOrders(orders []HistoryOrder) ([]order.HistoryOrder, error) {
	shortOrders := make([]order.HistoryOrder, 0, len(orders))

	for _, v := range orders {
		portShortOrder, err := ToPortHistoryOrder(v)
		if err != nil {
			return nil, fmt.Errorf("convert to short order: %w", err)
		}

		shortOrders = append(shortOrders, portShortOrder)
	}

	return shortOrders, nil
}

func ToPortHistoryOrder(dbOrder HistoryOrder) (order.HistoryOrder, error) {
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

func ToPortOrder(
	dbOrder *Order,
	dbOrderProducts []OrderProduct,
	dbTimeline []OrderTimeline,
) (*order.Order, error) {
	orderStatus, err := order.StatusFromString(dbOrder.Status)
	if err != nil {
		return nil, fmt.Errorf("status from string: %w", err)
	}

	portTimeline, err := toPortTimeline(dbTimeline)
	if err != nil {
		return nil, fmt.Errorf("timeline: %w", err)
	}

	return &order.Order{
		ID:               order.IDFromInt(dbOrder.ID),
		ChatID:           msginfo.ChatIDFromInt64(dbOrder.ChatID),
		Status:           orderStatus,
		VerificationCode: dbOrder.VerificationCode.String,
		CurrencyID:       currency.IDFromInt(dbOrder.CurrencyID),
		DailyPosition:    int(dbOrder.DailyPosition.Int32),
		TotalPrice:       dbOrder.TotalPrice,
		CreatedAt:        dbOrder.CreatedAt,
		UpdatedAt:        dbOrder.UpdatedAt,
		Timeline:         portTimeline,
		Products:         toPortCartProducts(dbOrderProducts),
	}, nil
}

func toPortTimeline(dbTimeline []OrderTimeline) ([]order.StatusTime, error) {
	sort.Slice(dbTimeline, func(i, j int) bool {
		return dbTimeline[i].UpdatedAt.Sub(dbTimeline[j].UpdatedAt) < 0
	})

	portTimeline := make([]order.StatusTime, 0, len(dbTimeline))

	for _, orderTime := range dbTimeline {
		status, err := order.StatusFromString(orderTime.Status)
		if err != nil {
			return nil, fmt.Errorf("timeline status from string: %w", err)
		}

		portTimeline = append(portTimeline, order.StatusTime{
			Status: status,
			Time:   orderTime.UpdatedAt,
		})
	}

	return portTimeline, nil
}
