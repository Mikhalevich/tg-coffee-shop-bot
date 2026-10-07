package setup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/app"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/manager/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/notificationsvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/ordersvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/nextpending"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/updatestatus"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func StartService(
	ctx context.Context,
	cfg config.Config,
) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	dbConn, driver, cleanup, err := MakePGXConnection(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make pgx connection: %w", err)
	}
	defer cleanup()

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, driver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		sender              = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		timeProvicer        = timeprovider.New()
		orderService        = ordersvc.New(
			transactionProvider,
			pgorder.New(driver, transactionProvider),
			timeProvicer,
		)
		notificationService = notificationsvc.New(
			pgoutbox.New(transactionProvider),
			sender,
		)
		nextPendingOrderUsecase = nextpending.New(
			transactionProvider,
			orderService,
			notificationService,
		)
		updateOrderStatusUsecase = updatestatus.New(
			transactionProvider,
			orderService,
			notificationService,
		)
	)

	if err := app.New(
		logger.FromContext(ctx),
		nextPendingOrderUsecase,
		updateOrderStatusUsecase,
	).Start(
		ctx,
		cfg.HTTPPort,
	); err != nil {
		return fmt.Errorf("start manager app: %w", err)
	}

	return nil
}

func MakePGXConnection(ctx context.Context, cfg config.Postgres) (*sql.DB, *driver.Pgx, func(), error) {
	if cfg.Connection == "" {
		return nil, nil, func() {}, nil
	}

	driver := driver.NewPgx()

	dbConn, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("ping: %w", err)
	}

	return dbConn, driver, func() {
		dbConn.Close()
	}, nil
}
