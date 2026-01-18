package orderpayment

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type MessageSender interface {
	SendMessage(ctx context.Context, msg messageprocessor.Message) error
	AnswerOrderPayment(
		ctx context.Context,
		paymentID string,
		ok bool,
		errorMsg string,
	) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type Transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error)
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)

	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)

	GetOrderPositionByStatus(ctx context.Context, id order.ID, statuses ...order.Status) (int, error)

	UpdateOrderByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		data order.UpdateOrderData,
		prevStatuses ...order.Status,
	) (*order.Order, error)

	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)

	IsNotFoundError(err error) bool
}

type StoreInfo interface {
	GetStoreByID(ctx context.Context, id store.ID) (*store.Store, error)
}

type TimeProvider interface {
	Now() time.Time
}

type VerificationCodeGenerator interface {
	Generate() string
}

type OrderPayment struct {
	storeID       store.ID
	sender        MessageSender
	escaper       MarkdownEscaper
	qrCode        port.QRCodeGenerator
	transactor    Transactor
	repository    Repository
	storeInfo     StoreInfo
	dailyPosition port.DailyPositionGenerator
	codeGenerator VerificationCodeGenerator
	timeProvider  TimeProvider
}

func New(
	storeID int,
	sender MessageSender,
	escaper MarkdownEscaper,
	qrCode port.QRCodeGenerator,
	transactor Transactor,
	repository Repository,
	storeInfo StoreInfo,
	dailyPosition port.DailyPositionGenerator,
	codeGenerator VerificationCodeGenerator,
	timeProvider TimeProvider,
) *OrderPayment {
	return &OrderPayment{
		storeID:       store.IDFromInt(storeID),
		sender:        sender,
		escaper:       escaper,
		qrCode:        qrCode,
		transactor:    transactor,
		repository:    repository,
		storeInfo:     storeInfo,
		dailyPosition: dailyPosition,
		codeGenerator: codeGenerator,
		timeProvider:  timeProvider,
	}
}
