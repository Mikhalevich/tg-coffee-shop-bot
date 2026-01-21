package setup

import (
	"context"
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
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
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

	pgDB, cleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}

	defer cleanup()

	var (
		sender          = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		msgProcessor    = messageprocessor.New(sender, sender, buttonRepository)
		timeProvider    = timeprovider.New()
		outboxProcessor = outboxprocessor.New(
			pgDB.Transactor(),
			pgDB,
			msgProcessor,
			timeProvider,
		)
	)

	app.New(outboxProcessor).Run(
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

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	driver := driver.NewPgx()

	dbConn, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, driver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		p                   = postgres.New(driver, transactionProvider)
	)

	return p, func() {
		dbConn.Close()
	}, nil
}
