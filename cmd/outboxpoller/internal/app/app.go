package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type Processor interface {
	ProcessMessage(ctx context.Context, batchSize int) error
	ProcessAnswerPayment(ctx context.Context, batchSize int) error
	ProcessInvoice(ctx context.Context, batchSize int) error
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
	invoiceCfg config.Worker,
) {
	var wgr sync.WaitGroup

	runWorkers(ctx, "message", messageCfg, &wgr,
		func(ctx context.Context) error {
			return a.processor.ProcessMessage(ctx, messageCfg.BatchSize)
		},
	)

	runWorkers(ctx, "answer_payment", answerPaymentCfg, &wgr,
		func(ctx context.Context) error {
			return a.processor.ProcessAnswerPayment(ctx, messageCfg.BatchSize)
		},
	)

	runWorkers(ctx, "invoice", invoiceCfg, &wgr,
		func(ctx context.Context) error {
			return a.processor.ProcessInvoice(ctx, messageCfg.BatchSize)
		},
	)

	wgr.Wait()
}

func runWorkers(
	ctx context.Context,
	workerName string,
	cfg config.Worker,
	wgr *sync.WaitGroup,
	processFn func(ctx context.Context) error,
) {
	for i := range cfg.Count {
		wgr.Go(func() {
			log := logger.FromContext(ctx).
				WithFields(logger.Fields{
					"worker_name":   workerName,
					"worker_number": i,
				})
			runPoller(
				logger.WithLogger(ctx, log),
				cfg.Interval,
				processFn,
			)
		})
	}
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
