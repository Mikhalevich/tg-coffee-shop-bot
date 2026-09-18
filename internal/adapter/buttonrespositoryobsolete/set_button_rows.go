package buttonrespositoryobsolete

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
)

func (r *ButtonRepository) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) error {
	cmds, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, row := range rows {
			for _, btn := range row {
				if err := r.storeButton(ctx, pipe, btn); err != nil {
					return fmt.Errorf("store button: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("pipelined: %w", err)
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return fmt.Errorf("pipeline cmd: %w", err)
		}
	}

	return nil
}
