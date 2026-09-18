package activeorder

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
	GetActiveOrder(ctx context.Context, chatID msginfo.ChatID) (order.Order, error)
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
	NoActiveOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	ViewOrder(
		ctx context.Context,
		chatID msginfo.ChatID,
		ord order.Order,
		products map[product.ProductID]product.Product,
		curr currency.Currency,
		position int,
	) error
}

type ActiveOrder struct {
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
) *ActiveOrder {
	return &ActiveOrder{
		orderService:        orderService,
		productService:      productService,
		currencyService:     currencyService,
		notificationService: notificationService,
	}
}

func (a *ActiveOrder) ViewActiveOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	activeOrder, err := a.orderService.GetActiveOrder(ctx, chatID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("get active order: %w", err)
		}

		if err := a.notificationService.NoActiveOrder(ctx, chatID); err != nil {
			return fmt.Errorf("no active order msg: %w", err)
		}

		return nil
	}

	productsInfo, err := a.productService.GetProductsByIDs(
		ctx,
		activeOrder.ProductIDs(),
		activeOrder.CurrencyID,
	)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	curr, err := a.currencyService.GetCurrencyByID(ctx, activeOrder.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := a.notificationService.ViewOrder(
		ctx,
		chatID,
		activeOrder,
		productsInfo,
		curr,
		a.orderQueuePosition(ctx, activeOrder),
	); err != nil {
		return fmt.Errorf("view order msg: %w", err)
	}

	return nil
}

func (a *ActiveOrder) orderQueuePosition(
	ctx context.Context,
	activeOrder order.Order,
) int {
	if !activeOrder.InQueue() {
		return 0
	}

	pos, err := a.orderService.GetActiveOrderPosition(
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
