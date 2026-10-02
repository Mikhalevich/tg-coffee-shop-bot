package payment

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type StoreService interface {
	GetStoreInfo(ctx context.Context) (store.StoreInfo, error)
}

type OrderService interface {
	GetOrderByID(
		ctx context.Context,
		id order.ID,
	) (order.Order, error)
	SetOrderPaymentInProgress(
		ctx context.Context,
		orderID order.ID,
		totalAmount int,
		operationTime time.Time,
	) error
	SetOrderPaymentConfirmed(
		ctx context.Context,
		chatID msginfo.ChatID,
		orderID order.ID,
		verificationCode string,
		dailyPosition int,
		operationTime time.Time,
	) (int, error)
}

type ProductsService interface {
	GetProductsByIDs(
		ctx context.Context,
		ids []product.ProductID,
		currencyID currency.ID,
	) (map[product.ProductID]product.Product, error)
}

type CurrencyService interface {
	GetCurrencyByID(
		ctx context.Context,
		id currency.ID,
	) (currency.Currency, error)
}

type PositionService interface {
	Position(ctx context.Context, t time.Time) (int, error)
}

type CodeGeneratorService interface {
	Generate() string
}

type QRCodeService interface {
	GeneratePNG(content string) ([]byte, error)
}

type TimeProvider interface {
	Now() time.Time
}

type NotificationService interface {
	PaymentStoreClosed(
		ctx context.Context,
		paymentID string,
		currentTime time.Time,
		nextWorkingTime time.Time,
	) error
	PaymentAnswer(
		ctx context.Context,
		paymentID string,
		success bool,
	) error
	PaymentConfirmed(
		ctx context.Context,
		chatID msginfo.ChatID,
		ord order.Order,
		curr currency.Currency,
		productsInfo map[product.ProductID]product.Product,
		queuePosition int,
		qrCodeImage []byte,
	) error
	SendOrderConfirmed(
		ctx context.Context,
		chatID msginfo.ChatID,
		ord order.Order,
		productsInfo map[product.ProductID]product.Product,
		curr currency.Currency,
		queuePosition int,
		imagePayload []byte,
	) error
}

type Payment struct {
	transactor           Transactor
	storeService         StoreService
	orderService         OrderService
	productsService      ProductsService
	currencyService      CurrencyService
	positionService      PositionService
	codeGeneratorService CodeGeneratorService
	qrCodeService        QRCodeService
	timeProvider         TimeProvider
	notificationService  NotificationService
}

func New(
	transactor Transactor,
	storeService StoreService,
	orderService OrderService,
	productsService ProductsService,
	currencyService CurrencyService,
	positionService PositionService,
	codeGeneratorService CodeGeneratorService,
	qrCodeService QRCodeService,
	timeProvider TimeProvider,
	notificationService NotificationService,
) *Payment {
	return &Payment{
		transactor:           transactor,
		storeService:         storeService,
		orderService:         orderService,
		productsService:      productsService,
		currencyService:      currencyService,
		positionService:      positionService,
		codeGeneratorService: codeGeneratorService,
		qrCodeService:        qrCodeService,
		timeProvider:         timeProvider,
		notificationService:  notificationService,
	}
}
