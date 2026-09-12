package cartorder

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (c *CartOrder) Cancel(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
) error {
	if err := c.cartService.Clear(ctx, chatID, cartID); err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return fmt.Errorf("clear cart: %w", err)
		}

		if err := c.notificationService.CartOrderUnavailable(ctx, chatID); err != nil {
			return fmt.Errorf("order unavailable msg: %w", err)
		}
	}

	return nil
}
