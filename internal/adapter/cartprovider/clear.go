package cartprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/redis/go-redis/v9"
)

func (c *CartProvider) Clear(ctx context.Context, chatID msginfo.ChatID, cartID cart.ID) error {
	if _, err := c.client.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		if err := c.client.Del(ctx, makeCartIDKey(chatID)).Err(); err != nil {
			return fmt.Errorf("cart id del: %w", err)
		}

		if err := c.client.Del(ctx, makeCartProductsKey(cartID.String())).Err(); err != nil {
			return fmt.Errorf("cart products del: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("pipelined: %w", err)
	}

	return nil
}
