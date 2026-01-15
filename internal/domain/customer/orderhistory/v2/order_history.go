package v2

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type Repository interface {
	HistoryOrdersByOffset(ctx context.Context, chatID msginfo.ChatID, offset, limit int) ([]order.HistoryOrder, error)
	HistoryOrdersCount(ctx context.Context, chatID msginfo.ChatID) (int, error)
}

type CurrencyProvider interface {
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
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
	repository       Repository
	currencyProvider CurrencyProvider
	sender           MessageSender
	pageSize         int
}

func New(
	repository Repository,
	currencyProvider CurrencyProvider,
	sender MessageSender,
	pageSize int,
) *OrderHistory {
	return &OrderHistory{
		repository:       repository,
		currencyProvider: currencyProvider,
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
