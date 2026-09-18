package ordersvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/orderbyid"
)

var (
	_ cartorder.OrderService = (*Service)(nil)
	_ orderbyid.OrderService = (*Service)(nil)
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

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
}

type Service struct {
	transactor Transactor
	repo       Repository
}

func New(
	transactor Transactor,
	repo Repository,
) *Service {
	return &Service{
		transactor: transactor,
		repo:       repo,
	}
}
