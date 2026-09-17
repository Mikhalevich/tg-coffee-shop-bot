package app

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/tracing"
)

type handlerFunc[I, O any] func(context.Context, *I) (*O, error)

func addRoute[I, O any](application *App, op huma.Operation, hf handlerFunc[I, O]) {
	huma.Register(
		application.humaAPI,
		op,
		makeHandlerWrapper(application, op.Path, hf),
	)
}

func makeHandlerWrapper[I, O any](application *App, pattern string, hndlrFn handlerFunc[I, O]) handlerFunc[I, O] {
	return func(ctx context.Context, input *I) (*O, error) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			log    = application.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog = logger.WithLogger(ctx, log)
		)

		output, err := hndlrFn(ctxLog, input)
		if err != nil {
			log.WithError(err).
				Error("http handler error")

			return output, supressSensitiveInfoFromError(err)
		}

		return output, nil
	}
}

func supressSensitiveInfoFromError(originErr error) error {
	if humaErr, ok := errors.AsType[*huma.ErrorModel](originErr); ok {
		humaErr.Errors = nil

		return humaErr
	}

	return huma.Error500InternalServerError("internal error")
}
