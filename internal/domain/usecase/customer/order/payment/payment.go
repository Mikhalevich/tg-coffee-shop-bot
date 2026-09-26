package payment

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type StoreService interface {
	GetStoreInfo(ctx context.Context) (store.StoreInfo, error)
}

type OrderService interface {
	SerOrderPaymentInProgress(
		ctx context.Context,
		orderID order.ID,
		totalAmount int,
		operationTime time.Time,
	) error
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
}

type Payment struct {
	transactor          Transactor
	storeService        StoreService
	orderService        OrderService
	notificationService NotificationService
}

func New(
	transactor Transactor,
	storeService StoreService,
	orderService OrderService,
	notificationService NotificationService,
) *Payment {
	return &Payment{
		transactor:          transactor,
		storeService:        storeService,
		orderService:        orderService,
		notificationService: notificationService,
	}
}
