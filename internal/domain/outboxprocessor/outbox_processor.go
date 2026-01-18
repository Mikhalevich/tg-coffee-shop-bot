package outboxprocessor

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
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
