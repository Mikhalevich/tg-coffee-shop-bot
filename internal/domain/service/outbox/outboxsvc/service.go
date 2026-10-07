package outboxsvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
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
		invoice order.Invoice,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type Service struct {
	transactor   Transactor
	repo         Repository
	sender       Sender
	timeProvider TimeProvider
}

func New(
	transactor Transactor,
	repo Repository,
	sender Sender,
	timeProvider TimeProvider,
) *Service {
	return &Service{
		transactor:   transactor,
		repo:         repo,
		sender:       sender,
		timeProvider: timeProvider,
	}
}
