package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/tracing"
)

func Start(
	ctx context.Context,
	botCfg config.Bot,
	cartUsecase tghandler.CartUsecase,
	viewActiveOrderUsecase tghandler.ViewActiveOrderUsecase,
	cancelOrderUsecase tghandler.CancelOrderUsecase,
	queueSizeUsecase tghandler.QueueSizeUsecase,
	historyOrderUsecase tghandler.HistoryOrderUsecase,
	orderPaymentUsecase tghandler.OrderPaymentUsecase,
	buttonProvider tghandler.ButtonProvider,
) error {
	var (
		botHandler = tghandler.New(
			cartUsecase,
			viewActiveOrderUsecase,
			cancelOrderUsecase,
			queueSizeUsecase,
			historyOrderUsecase,
			orderPaymentUsecase,
			buttonProvider,
		)
	)

	tbot, err := tgbot.New(
		botCfg.Token,
		tgbot.WithWebHookToken(botCfg.WebHookToken),
		tgbot.WithNewTracerFn(func() tgbot.Tracer {
			return newTracer(logger.FromContext(ctx))
		}),
	)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}

type tracer struct {
	log     logger.Logger
	endSpan func()
}

func newTracer(log logger.Logger) *tracer {
	return &tracer{
		log: log,
	}
}

func (t *tracer) Before(
	ctx context.Context,
	pattern string,
	msg tgbot.BotMessage,
) context.Context {
	ctx, span := tracing.StartSpanName(ctx, pattern)

	t.endSpan = func() {
		span.End()
	}

	log := t.log.WithContext(ctx).
		WithField("endpoint", pattern).
		WithField("bot_message", msg)

	return logger.WithLogger(ctx, log)
}

func (t *tracer) OnSuccess(ctx context.Context) {
}

func (t *tracer) OnError(ctx context.Context, err error) {
	logger.FromContext(ctx).
		WithError(err).
		Error("error while processing message")
}

func (t *tracer) After(ctx context.Context) {
	t.endSpan()
}
