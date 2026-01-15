package cartprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type StoreInfo interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}

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
	SendInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		title string,
		description string,
		ord *order.Order,
		productsInfo map[product.ProductID]product.Product,
		currency string,
		rows ...button.ButtonRow,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type CartProcessing struct {
	storeID      store.ID
	repository   port.CustomerCartRepository
	storeInfo    StoreInfo
	cart         port.Cart
	sender       MessageSender
	timeProvider TimeProvider
}

func New(
	storeID int,
	repository port.CustomerCartRepository,
	storeInfo StoreInfo,
	cart port.Cart,
	sender MessageSender,
	timeProvider TimeProvider,
) *CartProcessing {
	return &CartProcessing{
		storeID:      store.IDFromInt(storeID),
		repository:   repository,
		storeInfo:    storeInfo,
		cart:         cart,
		sender:       sender,
		timeProvider: timeProvider,
	}
}

func (cp *CartProcessing) sendPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
) {
	if err := cp.sender.SendTextPlain(ctx, chatID, text); err != nil {
		logger.FromContext(ctx).WithError(err).Error("send message")
	}
}

func (cp *CartProcessing) replyPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := cp.sender.ReplyTextPlain(
		ctx,
		chatID,
		replyMessageID,
		text,
		buttons...,
	); err != nil {
		logger.FromContext(ctx).WithError(err).Error("reply message")
	}
}

func (cp *CartProcessing) editPlainText(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	text string,
	buttons ...button.ButtonRow,
) {
	if err := cp.sender.EditMessage(ctx, chatID, messageID, text, buttons...); err != nil {
		logger.FromContext(ctx).WithError(err).Error("edit message")
	}
}
