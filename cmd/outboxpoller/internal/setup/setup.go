package setup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/app"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/outboxpoller/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/messagesvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/outbox/outboxsvc"
)

func StartPoller(
	ctx context.Context,
	cfg config.Config,
) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	buttonRepository, err := MakeRedisButtonRepository(ctx, cfg.ButtonRedis)
	if err != nil {
		return fmt.Errorf("make redis button repository: %w", err)
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
		timeProvider        = timeprovider.New()
		messageService      = messagesvc.New(
			sender,
			sender,
			buttonRepository,
		)
		outboxService = outboxsvc.New(
			transactionProvider,
			pgoutbox.New(transactionProvider),
			messageService,
			timeProvider,
		)
	)

	app.New(outboxService).Run(
		ctx,
		cfg.MessageWorker,
		cfg.AnswerPaymentWorker,
		cfg.InvoiceWorker,
	)

	return nil
}

func MakeRedisButtonRepository(
	ctx context.Context,
	cfg config.ButtonRedis,
) (*buttonrespository.ButtonRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return buttonrespository.New(rdb, cfg.TTL), nil
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
