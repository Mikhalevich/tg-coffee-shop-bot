package cartprovider

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
)

var _ cartprocessing.CartProvider = (*CartProvider)(nil)

type CartProvider struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *CartProvider {
	return &CartProvider{
		client: client,
		ttl:    ttl,
	}
}

func (c *CartProvider) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}
