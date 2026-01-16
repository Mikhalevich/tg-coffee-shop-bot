package orderprocessing

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type Transactor interface {
	Transaction(
		ctx context.Context,
		trxFn func(ctx context.Context) error,
	) error
}

type CustomerOutboxMessageSender interface {
	SendMessage(
		ctx context.Context,
		msg messageprocessor.Message,
	) error
}

type Repository interface {
	UpdateOrderStatusForMinID(
		ctx context.Context,
		operationTime time.Time,
		newStatus, prevStatus order.Status,
	) (*order.Order, error)
	UpdateOrderStatus(
		ctx context.Context,
		id order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) (*order.Order, error)
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type TimeProvider interface {
	Now() time.Time
}

type OrderProcessing struct {
	transactor     Transactor
	customerSender CustomerOutboxMessageSender
	repository     Repository
	escaper        MarkdownEscaper
	timeProvider   TimeProvider
}

func New(
	transactor Transactor,
	customerSender CustomerOutboxMessageSender,
	repository Repository,
	escaper MarkdownEscaper,
	timeProvider TimeProvider,
) *OrderProcessing {
	return &OrderProcessing{
		transactor:     transactor,
		customerSender: customerSender,
		repository:     repository,
		escaper:        escaper,
		timeProvider:   timeProvider,
	}
}
