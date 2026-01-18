package cartprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

const (
	emptyListPlaceholder = "empty_list_placeholder"
)

func (c *CartProvider) StartNewCart(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error) {
	var (
		cartIDKey       = makeCartIDKey(chatID)
		cartNewID       = generateID()
		cartProductsKey = makeCartProductsKey(cartNewID)
	)

	cartPrevID, err := c.activeCartID(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("get active cart id: %w", err)
	}

	if _, err := c.client.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		if cartPrevID.String() != "" {
			if err := pipeline.Del(ctx, makeCartProductsKey(cartPrevID.String())).Err(); err != nil {
				return fmt.Errorf("cart products del: %w", err)
			}
		}

		if err := pipeline.Set(ctx, cartIDKey, cartNewID, c.ttl).Err(); err != nil {
			return fmt.Errorf("set cart id: %w", err)
		}

		if err := pipeline.RPush(ctx, cartProductsKey, emptyListPlaceholder).Err(); err != nil {
			return fmt.Errorf("rpush: %w", err)
		}

		if err := pipeline.Expire(ctx, cartProductsKey, c.ttl).Err(); err != nil {
			return fmt.Errorf("products expire: %w", err)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("pipelined: %w", err)
	}

	return cart.IDFromString(cartNewID), nil
}

func (c *CartProvider) activeCartID(ctx context.Context, chatID msginfo.ChatID) (cart.ID, error) {
	cartID, err := c.client.Get(ctx, makeCartIDKey(chatID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", fmt.Errorf("get: %w", err)
	}

	return cart.IDFromString(cartID), nil
}

func generateID() string {
	return uuid.NewString()
}

func makeCartIDKey(chatID msginfo.ChatID) string {
	return fmt.Sprintf("cart:id:%d", chatID.Int64())
}

func makeCartProductsKey(id string) string {
	return fmt.Sprintf("cart:products:%s", id)
}

func isEmptyListPlaceholder(val string) bool {
	return val == emptyListPlaceholder
}
