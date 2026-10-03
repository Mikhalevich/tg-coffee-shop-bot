package history

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type OrderService interface {
	GetHistoryOrderPage(
		ctx context.Context,
		chatID msginfo.ChatID,
		page int,
		pageSize int,
	) ([]order.HistoryOrder, order.Page, error)
}

type CurrencyService interface {
	GetCurrencyByID(
		ctx context.Context,
		id currency.ID,
	) (currency.Currency, error)
}

type NotificationService interface {
	NoOrdersFound(ctx context.Context, chatID msginfo.ChatID) error
	ShowHistoryPage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		orders []order.HistoryOrder,
		pageInfo order.Page,
		curr currency.Currency,
	) error
}

type History struct {
	pageSize            int
	orderService        OrderService
	currencyService     CurrencyService
	notificationService NotificationService
}

func New(
	pageSize int,
	orderService OrderService,
	currencyService CurrencyService,
	notificationService NotificationService,
) *History {
	return &History{
		pageSize:            pageSize,
		orderService:        orderService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}
