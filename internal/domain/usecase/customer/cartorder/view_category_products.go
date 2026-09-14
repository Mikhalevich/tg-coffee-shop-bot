package cartorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (c *CartOrder) ViewCategoryProducts(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	categoryID product.CategoryID,
	currencyID currency.ID,
) error {
	cartProducts, err := c.cartService.GetProducts(ctx, cartID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("get cart products: %w", err)
		}

		if err := c.notificationService.CartOrderUnavailable(ctx, info.ChatID); err != nil {
			return fmt.Errorf("order unavailable msg: %w", err)
		}

		return nil
	}

	categoryProducts, err := c.productService.GetProductsByCategoryID(
		ctx,
		categoryID,
		currencyID,
	)

	if err != nil {
		return fmt.Errorf("get products by category id: %w", err)
	}

	curr, err := c.currencyService.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	if err := c.notificationService.ViewCategoryProducts(
		ctx,
		info.ChatID,
		info.MessageID,
		cartID,
		categoryID,
		categoryProducts,
		cartProducts,
		curr,
	); err != nil {
		return fmt.Errorf("view category products msg: %w", err)
	}

	return nil
}
