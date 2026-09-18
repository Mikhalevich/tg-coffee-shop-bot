package orderbyid

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, id order.ID) (order.Order, error)
	GetActiveOrderPosition(
		ctx context.Context,
		orderID order.ID,
	) (int, error)
}

type ProductService interface {
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

type NotificationService interface {
	InvalidOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	ViewOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
		ord order.Order,
		products map[product.ProductID]product.Product,
		curr currency.Currency,
		pos int,
	) error
}

type OrderByID struct {
	orderService        OrderService
	productService      ProductService
	currencyService     CurrencyService
	notificationService NotificationService
}

func New(
	orderService OrderService,
	productService ProductService,
	currencyService CurrencyService,
	notificationService NotificationService,
) *OrderByID {
	return &OrderByID{
		orderService:        orderService,
		productService:      productService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}

func (o *OrderByID) GetOrderByID(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderID order.ID,
) error {
	ord, err := o.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("get order: %w", err)
		}

		if err := o.notificationService.InvalidOrder(ctx, chatID); err != nil {
			return fmt.Errorf("send invalid order: %w", err)
		}

		return nil
	}

	if !ord.IsSameChat(chatID) {
		if err := o.notificationService.InvalidOrder(ctx, chatID); err != nil {
			return fmt.Errorf("send invalid order: %w", err)
		}

		return nil
	}

	productsInfo, err := o.productService.GetProductsByIDs(
		ctx,
		ord.ProductIDs(),
		ord.CurrencyID,
	)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	curr, err := o.currencyService.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := o.notificationService.ViewOrder(
		ctx,
		chatID,
		ord,
		productsInfo,
		curr,
		o.orderQueuePosition(ctx, ord),
	); err != nil {
		return fmt.Errorf("view order msg: %w", err)
	}

	return nil
}

func (o *OrderByID) orderQueuePosition(
	ctx context.Context,
	activeOrder order.Order,
) int {
	if !activeOrder.InQueue() {
		return 0
	}

	pos, err := o.orderService.GetActiveOrderPosition(
		ctx,
		activeOrder.ID,
	)

	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			return 0
		}

		logger.FromContext(ctx).
			WithError(err).
			Error("failed to get order position")

		return 0
	}

	return pos
}
