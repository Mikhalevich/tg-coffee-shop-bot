package cartprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type CreateOrderInput struct {
	ChatID              msginfo.ChatID
	Status              order.Status
	StatusOperationTime time.Time
	VerificationCode    string
	TotalPrice          int
	Products            []order.OrderedProduct
	CurrencyID          currency.ID
}

type Transactor interface {
	Transaction(
		ctx context.Context,
		trxFn func(ctx context.Context) error,
	) error
}

type Repository interface {
	CreateOrder(ctx context.Context, coi CreateOrderInput) (*order.Order, error)
	GetCategories(ctx context.Context) ([]product.Category, error)
	GetProductsByCategoryID(
		ctx context.Context,
		categoryID product.CategoryID,
		currencyID currency.ID,
	) ([]product.Product, error)
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
	IsAlreadyExistsError(err error) bool
}

type StoreInfo interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}

type CartProvider interface {
	StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error)
	GetProducts(ctx context.Context, id cart.ID) ([]cart.CartProduct, error)
	AddProduct(ctx context.Context, id cart.ID, p cart.CartProduct) error
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
	IsNotFoundError(err error) bool
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
}

type InvoiceSender interface {
	SendInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		title string,
		ord *order.Order,
		productsInfo map[product.ProductID]product.Product,
		curr *currency.Currency,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type CartProcessing struct {
	storeID       store.ID
	transactor    Transactor
	repository    Repository
	storeInfo     StoreInfo
	cart          CartProvider
	sender        MessageSender
	invoiceSender InvoiceSender
	timeProvider  TimeProvider
}

func New(
	storeID int,
	transactor Transactor,
	repository Repository,
	storeInfo StoreInfo,
	cart CartProvider,
	sender MessageSender,
	invoiceSender InvoiceSender,
	timeProvider TimeProvider,
) *CartProcessing {
	return &CartProcessing{
		storeID:       store.IDFromInt(storeID),
		transactor:    transactor,
		repository:    repository,
		storeInfo:     storeInfo,
		cart:          cart,
		sender:        sender,
		invoiceSender: invoiceSender,
		timeProvider:  timeProvider,
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
