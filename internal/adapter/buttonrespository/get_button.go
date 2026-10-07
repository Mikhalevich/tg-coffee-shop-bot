package buttonrespository

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (r *ButtonRepository) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	key, btnNum := parseButtonID(id)
	if btnNum == "" {
		btn, err := r.singleButton(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("single button: %w", err)
		}

		return btn, nil
	}

	btn, err := r.hmapButton(ctx, key, btnNum)
	if err != nil {
		return nil, fmt.Errorf("hmap button: %w", err)
	}

	return btn, nil
}

func (r *ButtonRepository) singleButton(ctx context.Context, key string) (*button.Button, error) {
	b, err := r.client.GetDel(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	btn, err := decodeButton(b)
	if err != nil {
		return nil, fmt.Errorf("decode button: %w", err)
	}

	return btn, nil
}

func (r *ButtonRepository) hmapButton(ctx context.Context, key, field string) (*button.Button, error) {
	b, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, fmt.Errorf("hget: %w", err)
	}

	btn, err := decodeButton([]byte(b))
	if err != nil {
		return nil, fmt.Errorf("decode button: %w", err)
	}

	return btn, nil
}
