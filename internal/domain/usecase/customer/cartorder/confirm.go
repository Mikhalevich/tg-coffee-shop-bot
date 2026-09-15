package cartorder

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

//nolint:cyclop,funlen
func (c *CartOrder) Confirm(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	currencyID currency.ID,
) error {
	storeInfo, err := c.storeService.GetStoreInfo(ctx)
	if err != nil {
		return fmt.Errorf("get store info: %w", err)
	}

	if !storeInfo.IsActive {
		if err := c.notificationService.SendStoreClosed(
			ctx,
			info.ChatID,
			storeInfo.CurrentTime,
			storeInfo.NextWorkingTime,
		); err != nil {
			return fmt.Errorf("store info closed msg: %w", err)
		}

		return nil
	}

	curr, err := c.currencyService.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	orderedProducts, productsInfo, err := c.orderedProductsFromCart(
		ctx,
		cartID,
		currencyID,
	)

	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("ordered products from cart: %w", err)
		}

		if err := c.notificationService.CartOrderUnavailable(ctx, info.ChatID); err != nil {
			return fmt.Errorf("cart order unavailable msg: %w", err)
		}

		return nil
	}

	if len(orderedProducts) == 0 {
		if err := c.notificationService.NoProductsForOrder(ctx, info.ChatID); err != nil {
			return fmt.Errorf("no products for order msg: %w", err)
		}

		return nil
	}

	if err := c.createOrder(
		ctx,
		info.ChatID,
		orderedProducts,
		productsInfo,
		*curr,
		storeInfo.CurrentTime,
	); err != nil {
		if !perror.IsType(err, perror.TypeAlreadyExists) {
			return fmt.Errorf("create order: %w", err)
		}

		if err := c.notificationService.OrderAlreadyExists(ctx, info.ChatID); err != nil {
			return fmt.Errorf("order already exists msg: %w", err)
		}

		return nil
	}

	if err := c.cartService.Clear(ctx, info.ChatID, cartID); err != nil {
		return fmt.Errorf("cart clear: %w", err)
	}

	return nil
}

func (c *CartOrder) createOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
	curr currency.Currency,
	now time.Time,
) error {
	if err := c.transactor.Transaction(ctx, func(ctx context.Context) error {
		createdOrder, err := c.orderService.CreateOrder(
			ctx,
			c.makeCreateOrderInfo(chatID, orderedProducts, curr.ID, now),
		)
		if err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		if err := c.notificationService.SendInvoice(
			ctx,
			chatID,
			createdOrder,
			productsInfo,
			curr,
		); err != nil {
			return fmt.Errorf("send order invoice msg: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (c *CartOrder) makeCreateOrderInfo(
	chatID msginfo.ChatID,
	orderedProducts []order.OrderedProduct,
	currencyID currency.ID,
	now time.Time,
) order.CreateOrderInfo {
	totalPrice := 0
	for _, v := range orderedProducts {
		totalPrice += v.Count * v.Price
	}

	return order.CreateOrderInfo{
		ChatID:              chatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: now,
		VerificationCode:    "",
		TotalPrice:          totalPrice,
		Products:            orderedProducts,
		CurrencyID:          currencyID,
	}
}
