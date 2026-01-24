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

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/app"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/cartprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/dailypositiongenerator"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/qrcodegenerator"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/orderhistoryid"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/orderhistoryoffset"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/verificationcodegenerator"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory"
	orderhistoryv2 "github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderhistory/v2"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderpayment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
)

func StartBot(ctx context.Context, cfg config.Config) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	dbConn, driver, cleanup, err := MakePGXConnection(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make pgx connection: %w", err)
	}
	defer cleanup()

	cartRedis, err := MakeRedisCart(ctx, cfg.CartRedis)
	if err != nil {
		return fmt.Errorf("make redis cart: %w", err)
	}

	dailyPosition, err := MakeRedisDailyPositionGenerator(ctx, cfg.DailyPositionRedis)
	if err != nil {
		return fmt.Errorf("make redis daily position generator: %w", err)
	}

	buttonRepository, err := MakeRedisButtonRepository(ctx, cfg.ButtonRedis)
	if err != nil {
		return fmt.Errorf("make redis button repository: %w", err)
	}

	var (
		sqlxDBConn          = sqlx.NewDb(dbConn, driver.Name())
		transactionProvider = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		pgDB                = postgres.New(driver, transactionProvider)
		pgOrderHistoryID    = orderhistoryid.New(dbConn, driver)
		pgOrderHistoryPage  = orderhistoryoffset.New(dbConn, driver)
		sender              = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		msgProcessor        = messageprocessor.New(sender, sender, buttonRepository)
		qrGenerator         = qrcodegenerator.New()
		timeProvider        = timeprovider.New()
		cartProcessor       = cartprocessing.New(cfg.StoreID, transactionProvider,
			pgDB, pgDB, cartRedis, msgProcessor, pgDB, timeProvider)
		actionProcessor    = orderaction.New(msgProcessor, pgDB, timeProvider)
		historyProcessor   = orderhistory.New(pgDB, pgOrderHistoryID, msgProcessor, cfg.OrderHistory.PageSize)
		historyProcessorV2 = orderhistoryv2.New(pgOrderHistoryPage, pgDB, msgProcessor, cfg.OrderHistory.PageSize)
		paymentProcessor   = orderpayment.New(cfg.StoreID, pgDB, msgProcessor, qrGenerator,
			pgDB.Transactor(), pgDB, pgDB, dailyPosition, verificationcodegenerator.New(), timeProvider)
	)

	if err := app.Start(
		ctx,
		cfg.Bot.Token,
		cartProcessor,
		actionProcessor,
		historyProcessor,
		historyProcessorV2,
		paymentProcessor,
		msgProcessor,
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

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

func MakeRedisCart(ctx context.Context, cfg config.CartRedis) (cartprocessing.CartProvider, error) {
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

	return cartprovider.New(rdb, cfg.TTL), nil
}

func MakeRedisDailyPositionGenerator(
	ctx context.Context,
	cfg config.DailyPositionRedis,
) (orderpayment.DailyPositionGenerator, error) {
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

	return dailypositiongenerator.New(rdb, cfg.TTL), nil
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
