package orderaction

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type MessageSender interface {
	SendTextPlain(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		rows ...button.ButtonRow,
	) error
	ReplyTextPlain(
		ctx context.Context,
		chatID msginfo.ChatID,
		replyMessageID msginfo.MessageID,
		text string,
		rows ...button.ButtonRow,
	) error
	ReplyTextMarkdown(
		ctx context.Context,
		chatID msginfo.ChatID,
		replyMessageID msginfo.MessageID,
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
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
	EscapeMarkdown(s string) string
}

type Repository interface {
	GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error)
	GetOrderByChatIDAndStatus(ctx context.Context, id msginfo.ChatID, statuses ...order.Status) (*order.Order, error)
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)

	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)

	GetOrderPositionByStatus(ctx context.Context, id order.ID, statuses ...order.Status) (int, error)
	GetOrdersCountByStatus(ctx context.Context, statuses ...order.Status) (int, error)

	UpdateOrderByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		data order.UpdateOrderData,
		prevStatuses ...order.Status,
	) (*order.Order, error)

	UpdateOrderStatusByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)

	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
}

type TimeProvider interface {
	Now() time.Time
}

type OrderAction struct {
	sender       MessageSender
	repository   Repository
	timeProvider TimeProvider
}

func New(
	sender MessageSender,
	repository Repository,
	timeProvider TimeProvider,
) *OrderAction {
	return &OrderAction{
		sender:       sender,
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (o *OrderAction) deleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) {
	if err := o.sender.DeleteMessage(ctx, chatID, messageID); err != nil {
		logger.FromContext(ctx).WithError(err).Error("delete message")
	}
}

func (o *OrderAction) sendPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
) {
	if err := o.sender.SendTextPlain(ctx, chatID, text); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message")
	}
}

func (o *OrderAction) replyPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.ReplyTextPlain(
		ctx,
		chatID,
		replyMessageID,
		text,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("reply message")
	}
}

func (o *OrderAction) replyMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := o.sender.ReplyTextMarkdown(
		ctx,
		chatID,
		replyMessageID,
		text,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("reply message")
	}
}

func (o *OrderAction) editPlainText(
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
