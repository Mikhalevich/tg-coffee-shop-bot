package buttonrespository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (r *ButtonRepository) SetButton(ctx context.Context, btn button.Button) error {
	if err := r.storeButton(ctx, r.client, btn); err != nil {
		return fmt.Errorf("store button: %w", err)
	}

	return nil
}

func (r *ButtonRepository) storeButton(
	ctx context.Context,
	cmd redis.StringCmdable,
	btn button.Button,
) error {
	encodedButton, err := encodeButton(btn)
	if err != nil {
		return fmt.Errorf("encode button: %w", err)
	}

	if err := cmd.Set(ctx, btn.ID.String(), encodedButton, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
