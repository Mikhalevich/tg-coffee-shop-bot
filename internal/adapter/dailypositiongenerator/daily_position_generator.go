package dailypositiongenerator

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderpayment"
)

var _ orderpayment.DailyPositionGenerator = (*RedisDeilyPositionGenerator)(nil)

type RedisDeilyPositionGenerator struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *RedisDeilyPositionGenerator {
	return &RedisDeilyPositionGenerator{
		client: client,
		ttl:    ttl,
	}
}

func (r *RedisDeilyPositionGenerator) Position(ctx context.Context, t time.Time) (int, error) {
	dailyKey := makeKey(t)

	pos, err := r.client.Incr(ctx, dailyKey).Result()
	if err != nil {
		return 0, fmt.Errorf("incr: %w", err)
	}

	if pos == 1 {
		if err := r.client.Expire(ctx, dailyKey, r.ttl).Err(); err != nil {
			return 0, fmt.Errorf("expire: %w", err)
		}
	}

	return int(pos), nil
}

func makeKey(t time.Time) string {
	y, m, d := t.Date()

	return fmt.Sprintf("dailyposition:%d_%d_%d", y, m, d)
}
