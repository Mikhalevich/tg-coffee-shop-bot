package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/null"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
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

func (o Order) ToDom(
	dbOrderProducts []OrderProduct,
	dbTimeline []OrderTimeline,
) (*order.Order, error) {
	orderStatus, err := order.StatusFromString(o.Status)
	if err != nil {
		return nil, fmt.Errorf("status from string: %w", err)
	}

	portTimeline, err := toDomTimeline(dbTimeline)
	if err != nil {
		return nil, fmt.Errorf("timeline: %w", err)
	}

	return &order.Order{
		ID:               order.IDFromInt(o.ID),
		ChatID:           msginfo.ChatIDFromInt64(o.ChatID),
		Status:           orderStatus,
		VerificationCode: o.VerificationCode.String,
		CurrencyID:       currency.IDFromInt(o.CurrencyID),
		DailyPosition:    int(o.DailyPosition.Int32),
		TotalPrice:       o.TotalPrice,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
		Timeline:         portTimeline,
		Products:         toDomOrderedProducts(dbOrderProducts),
	}, nil
}

func ToDBOrder(domOrder order.Order) Order {
	return Order{
		ID:               domOrder.ID.Int(),
		ChatID:           domOrder.ChatID.Int64(),
		Status:           domOrder.Status.String(),
		VerificationCode: null.NullString(domOrder.VerificationCode),
		CurrencyID:       domOrder.CurrencyID.Int(),
		//nolint:gosec
		DailyPosition: null.NullIntPositive(int32(domOrder.DailyPosition)),
		TotalPrice:    domOrder.TotalPrice,
		CreatedAt:     domOrder.CreatedAt,
		UpdatedAt:     domOrder.UpdatedAt,
	}
}
