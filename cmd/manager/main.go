package main

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/setup"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/application"
)

func main() {
	var cfg config.Config
	application.Run(&cfg, func(ctx context.Context) error {
		if err := setup.StartService(ctx, cfg); err != nil {
			return fmt.Errorf("start service: %w", err)
		}

		return nil
	})
}
