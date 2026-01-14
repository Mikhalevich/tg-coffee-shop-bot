package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/setup"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/tracing"
)

func main() {
	var cfg config.Config
	if err := infra.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	log, err := infra.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.StdLogger().WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	log = log.WithField("service_name", cfg.Tracing.ServiceName)

	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		log.WithError(err).Error("failed to setup tracer")
		os.Exit(1)
	}

	if err := infra.RunSignalInterruptionFunc(func(ctx context.Context) error {
		log.Info("service starting")
		defer log.Info("service stopped")

		if err := setup.StartPoller(
			logger.WithLogger(ctx, log),
			cfg,
		); err != nil {
			return fmt.Errorf("start poller: %w", err)
		}

		return nil
	}); err != nil {
		log.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}
