package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type Processor interface {
	Process(ctx context.Context, batchSize int) error
	ProcessAnswerPayment(ctx context.Context, batchSize int) error
}

type App struct {
	processor Processor
}

func New(processor Processor) *App {
	return &App{
		processor: processor,
	}
}

func (a *App) Run(
	ctx context.Context,
	messageCfg config.Worker,
	answerPaymentCfg config.Worker,
) {
	var wgr sync.WaitGroup

	for i := range messageCfg.Count {
		wgr.Go(func() {
			log := logger.FromContext(ctx).
				WithFields(logger.Fields{
					"worker_name":   "message",
					"worker_number": i,
				})
			runPoller(
				logger.WithLogger(ctx, log),
				messageCfg.Interval,
				func(ctx context.Context) error {
					return a.processor.Process(ctx, messageCfg.BatchSize)
				},
			)
		})
	}

	for i := range answerPaymentCfg.Count {
		wgr.Go(func() {
			log := logger.FromContext(ctx).
				WithFields(logger.Fields{
					"worker_name":   "answer_payment",
					"worker_number": i,
				})
			runPoller(
				logger.WithLogger(ctx, log),
				answerPaymentCfg.Interval,
				func(ctx context.Context) error {
					return a.processor.ProcessAnswerPayment(ctx, messageCfg.BatchSize)
				},
			)
		})
	}

	wgr.Wait()
}

func runPoller(
	ctx context.Context,
	interval time.Duration,
	processFn func(ctx context.Context) error,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := processFn(ctx); err != nil {
				logger.FromContext(ctx).
					WithError(err).
					Error("process error")
			}

		case <-ctx.Done():
			return
		}
	}
}
