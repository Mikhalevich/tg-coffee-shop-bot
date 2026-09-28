package outboxsvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Repository interface {
	SelectForDispatchMessages(
		ctx context.Context,
		limit int,
	) ([]outboxmsg.Message, error)

	SetDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error

	SelectForDispatchAnswerPayment(
		ctx context.Context,
		limit int,
	) ([]outboxmsg.AnswerPayment, error)

	SetAnswerPaymentDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error

	SelectForDispatchInvoice(
		ctx context.Context,
		limit int,
	) ([]outboxmsg.Invoice, error)

	SetInvoiceDispatched(
		ctx context.Context,
		ids []int,
		dispatchedAt time.Time,
	) error
}

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg msginfo.Message,
	) error

	AnswerPayment(
		ctx context.Context,
		paymentID string,
		ok bool,
		errorMsg string,
	) error

	SendInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		title string,
		ord order.Order,
		productsInfo map[product.ProductID]product.Product,
		curr currency.Currency,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type OrderService interface {
	GetOrderByID(
		ctx context.Context,
		id order.ID,
	) (order.Order, error)
}

type CurrencyService interface {
	GetCurrencyByID(
		ctx context.Context,
		id currency.ID,
	) (currency.Currency, error)
}

type ProductsService interface {
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
}

type Service struct {
	transactor      Transactor
	repo            Repository
	sender          Sender
	timeProvider    TimeProvider
	orderService    OrderService
	currencyService CurrencyService
	productsService ProductsService
}

func New(
	transactor Transactor,
	repo Repository,
	sender Sender,
	timeProvider TimeProvider,
	orderService OrderService,
	currencyService CurrencyService,
	productsService ProductsService,
) *Service {
	return &Service{
		transactor:      transactor,
		repo:            repo,
		sender:          sender,
		timeProvider:    timeProvider,
		orderService:    orderService,
		currencyService: currencyService,
		productsService: productsService,
	}
}
