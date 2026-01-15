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
}

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg messageprocessor.Message,
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
