package orderhistory

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type CurrencyProvider interface {
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
}

type Repository interface {
	HistoryOrdersCount(ctx context.Context, chatID msginfo.ChatID) (int, error)
	HistoryOrdersFirst(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersLast(ctx context.Context, chatID msginfo.ChatID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersBeforeID(ctx context.Context, chatID msginfo.ChatID, id order.ID, size int) ([]order.HistoryOrder, error)
	HistoryOrdersAfterID(ctx context.Context, chatID msginfo.ChatID, id order.ID, size int) ([]order.HistoryOrder, error)
}

type MessageSender interface {
	SendTextPlain(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		rows ...button.ButtonRow,
	) error
	EditMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		text string,
		rows ...button.ButtonRow,
	) error
}

type OrderHistory struct {
	currencyProvider CurrencyProvider
	repository       Repository
	sender           MessageSender
	pageSize         int
}

func New(
	currencyProvider CurrencyProvider,
	repository Repository,
	sender MessageSender,
	pageSize int,
) *OrderHistory {
	return &OrderHistory{
		currencyProvider: currencyProvider,
		repository:       repository,
		sender:           sender,
		pageSize:         pageSize,
	}
}

func (o *OrderHistory) sendPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.SendTextPlain(ctx, chatID, text, buttons...); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message plain")
	}
}

func (o *OrderHistory) editPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.EditMessage(ctx, chatID, messageID, text, buttons...); err != nil {
		logger.FromContext(ctx).WithError(err).Error("edit message")
	}
}
