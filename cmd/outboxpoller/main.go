package main

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/setup"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/application"
)

func main() {
	var cfg config.Config

	application.Run(&cfg, func(ctx context.Context) error {
		if err := setup.StartPoller(ctx, cfg); err != nil {
			return fmt.Errorf("start poller: %w", err)
		}

		return nil
	})
}
