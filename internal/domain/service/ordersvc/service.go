package ordersvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/history"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/orderbyid"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/ordercancel"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/payment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/queuesize"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/nextpending"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/updatestatus"
)

var (
	_ cartorder.OrderService    = (*Service)(nil)
	_ orderbyid.OrderService    = (*Service)(nil)
	_ queuesize.OrderService    = (*Service)(nil)
	_ ordercancel.OrderService  = (*Service)(nil)
	_ payment.OrderService      = (*Service)(nil)
	_ history.OrderService      = (*Service)(nil)
	_ nextpending.OrderService  = (*Service)(nil)
	_ updatestatus.OrderService = (*Service)(nil)
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

//nolint:interfacebloat
type Repository interface {
	InsertOrder(
		ctx context.Context,
		ord order.Order,
	) (order.ID, error)
	InsertProductsOrder(
		ctx context.Context,
		orderID order.ID,
		products []order.OrderedProduct,
	) error
	InsertOrderTimeline(
		ctx context.Context,
		orderID order.ID,
		status order.Status,
		createdAt time.Time,
	) error
	GetOrderByID(ctx context.Context, id order.ID) (order.Order, error)
	GetOrderByChatIDAndStatus(
		ctx context.Context,
		id msginfo.ChatID,
		statuses ...order.Status,
	) (order.Order, error)
	GetOrderPositionByStatus(
		ctx context.Context,
		orderID order.ID,
		statuses ...order.Status,
	) (int, error)
	GetOrdersCountByStatus(
		ctx context.Context,
		statuses ...order.Status,
	) (int, error)
	UpdateOrderStatusByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) error
	UpdateOrderStatus(
		ctx context.Context,
		orderID order.ID,
		operationTime time.Time,
		newStatus order.Status,
		prevStatuses ...order.Status,
	) error
	UpdateOrderByChatAndID(
		ctx context.Context,
		orderID order.ID,
		chatID msginfo.ChatID,
		data order.UpdateOrderData,
		prevStatuses ...order.Status,
	) error
	HistoryOrdersByOffset(
		ctx context.Context,
		chatID msginfo.ChatID,
		offset int,
		limit int,
	) ([]order.HistoryOrder, error)
	HistoryOrdersCount(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (int, error)
	UpdateOrderStatusForMinID(
		ctx context.Context,
		operationTime time.Time,
		newStatus, prevStatus order.Status,
	) (order.ID, error)
}

type TimeProvider interface {
	Now() time.Time
}

type Service struct {
	transactor   Transactor
	repo         Repository
	timeProvider TimeProvider
}

func New(
	transactor Transactor,
	repo Repository,
	timeProvider TimeProvider,
) *Service {
	return &Service{
		transactor:   transactor,
		repo:         repo,
		timeProvider: timeProvider,
	}
}
