package outboxprocessor

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type OutboxMessage struct {
	messageprocessor.Message

	ID int
}

type OutboxAnswerPayment struct {
	ID        int
	PaymentID string
	OK        bool
	ErrorMsg  string
}

type OutboxInvoice struct {
	ID      int
	ChatID  msginfo.ChatID
	Text    string
	OrderID order.ID
}

type Transactor interface {
	Transaction(
		ctx context.Context,
		trxFn func(ctx context.Context) error,
	) error
}

type Repository interface {
	OutboxSelectForDispatchMessages(
		ctx context.Context,
		limit int,
	) ([]OutboxMessage, error)

	OutboxSetDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error

	OutboxSelectForDispatchAnswerPayment(
		ctx context.Context,
		limit int,
	) ([]OutboxAnswerPayment, error)

	OutboxSetAnswerPaymentDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error

	OutboxSelectForDispatchInvoice(
		ctx context.Context,
		limit int,
	) ([]OutboxInvoice, error)

	OutboxSetInvoiceDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error

	GetOrderByID(ctx context.Context, id order.ID) (*order.Order, error)
	GetCurrencyByID(ctx context.Context, id currency.ID) (*currency.Currency, error)
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
}

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg messageprocessor.Message,
	) error

	AnswerOrderPayment(
		ctx context.Context,
		paymentID string,
		ok bool,
		errorMsg string,
	) error

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

type OutboxProcessor struct {
	transactor   Transactor
	repository   Repository
	sender       Sender
	timeProvider TimeProvider
}

func New(
	transactor Transactor,
	repository Repository,
	sender Sender,
	timeProvider TimeProvider,
) *OutboxProcessor {
	return &OutboxProcessor{
		transactor:   transactor,
		repository:   repository,
		sender:       sender,
		timeProvider: timeProvider,
	}
}
