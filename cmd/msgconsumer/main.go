package main

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/setup"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/application"
)

func main() {
	var cfg config.Config
	application.Run(&cfg, func(ctx context.Context) error {
		if err := setup.StartConsumer(ctx, cfg); err != nil {
			return fmt.Errorf("start consumer: %w", err)
		}

		return nil
	})
}
