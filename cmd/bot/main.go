package main

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/setup"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/application"
)

func main() {
	var cfg config.Config
	application.Run(&cfg, func(ctx context.Context) error {
		if err := setup.StartBot(ctx, cfg); err != nil {
			return fmt.Errorf("start bot: %w", err)
		}

		return nil
	})
}
