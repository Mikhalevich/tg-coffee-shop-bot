package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func Start(
	ctx context.Context,
	token string,
	cartProcessor tghandler.CartProcessor,
	actionProcessor tghandler.OrderActionProcessor,
	historyProcessor tghandler.OrderHistoryProcessor,
	historyProcessorV2 tghandler.OrderHistoryProcessorV2,
	paymentProcessor tghandler.OrderPaymentProcessor,
	buttonProvider tghandler.ButtonProvider,
) error {
	var (
		botHandler = tghandler.New(
			cartProcessor,
			actionProcessor,
			historyProcessor,
			historyProcessorV2,
			paymentProcessor,
			buttonProvider,
		)
	)

	tbot, err := tgbot.New(token, logger.FromContext(ctx))
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
