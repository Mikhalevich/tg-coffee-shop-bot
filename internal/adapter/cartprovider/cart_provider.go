package cartprovider

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/cartsvc"
)

var (
	_ cartsvc.Repository = (*CartProvider)(nil)
)

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
