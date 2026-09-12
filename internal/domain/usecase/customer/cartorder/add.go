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

func (c *CartOrder) Add(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
	categoryID product.CategoryID,
	productID product.ProductID,
	currencyID currency.ID,
) error {
	if err := c.cartService.AddProduct(ctx, cartID, cart.CartProduct{
		ProductID:  productID,
		CategoryID: categoryID,
		Count:      1,
	}); err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("add cart product: %w", err)
		}

		if err := c.notificationService.CartOrderUnavailable(ctx, chatID); err != nil {
			return fmt.Errorf("order unavailable msg: %w", err)
		}
	}

	return nil
}
