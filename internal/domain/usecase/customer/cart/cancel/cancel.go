package cancel

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

type CartService interface {
	Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error
}

type NotificationService interface {
	CartOrderUnavailable(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
}

type CancelCart struct {
	cartService         CartService
	notificationService NotificationService
}

func New(
	cartService CartService,
	notificationService NotificationService,
) *CancelCart {
	return &CancelCart{
		cartService:         cartService,
		notificationService: notificationService,
	}
}

func (c *CancelCart) Cancel(
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
