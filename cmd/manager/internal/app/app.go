package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/app/handler"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

const (
	readTimeout      = time.Second * 10
	writeTimeout     = time.Second * 10
	shoutdownTimeout = time.Second * 30
)

type App struct {
	mux     *http.ServeMux
	humaAPI huma.API
	logger  logger.Logger
}

func New(
	logger logger.Logger,
	getNextOrderUsecase handler.GetNextOrderUsecase,
	updateOrderStatusUsecase handler.UpdateOrderStatusUsecase,
) *App {
	var (
		mux     = http.NewServeMux()
		humaAPI = humago.New(mux, huma.DefaultConfig("Bonus points", "1.0.0"))
	)

	httpManager := &App{
		mux:     mux,
		humaAPI: humaAPI,
		logger:  logger,
	}

	httpManager.routes(
		handler.New(
			getNextOrderUsecase,
			updateOrderStatusUsecase,
		),
	)

	return httpManager
}

func (application *App) Start(
	ctx context.Context,
	port int,
) error {
	var (
		srv = &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      application.mux,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		}

		srvErrCh = make(chan error)
	)

	defer close(srvErrCh)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrCh <- err
		}
	}()

	select {
	case err := <-srvErrCh:
		return fmt.Errorf("listen and serve: %w", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shoutdownTimeout)
	defer cancel()

	//nolint:contextcheck
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
