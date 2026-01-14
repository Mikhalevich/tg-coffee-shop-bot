package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type Processor interface {
	Process(ctx context.Context, batchSize int) error
}

type App struct {
	processor Processor
}

func New(processor Processor) *App {
	return &App{
		processor: processor,
	}
}

func (a *App) Run(ctx context.Context, count int, interval time.Duration, batchSize int) {
	var wgr sync.WaitGroup

	for i := range count {
		wgr.Go(func() {
			log := logger.FromContext(ctx).WithField("worker_number", i)
			a.runPoller(logger.WithLogger(ctx, log), interval, batchSize)
		})
	}

	wgr.Wait()
}

func (a *App) runPoller(ctx context.Context, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.processor.Process(ctx, batchSize); err != nil {
				logger.FromContext(ctx).
					WithError(err).
					Error("process error")
			}

		case <-ctx.Done():
			return
		}
	}
}
