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
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgcurrency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgoutbox"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgproduct"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/pgstore"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/verificationcodegenerator"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/cartprocessing"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/customer/orderpayment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/store"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/cartsvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/currencysvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/notificationsvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/ordersvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/productsvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/storesvc"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/history"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/payment"
)

//nolint:funlen
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
		sender              = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		timeProvider        = timeprovider.New()
		storeService        = storesvc.New(
			store.IDFromInt(cfg.StoreID),
			pgstore.New(transactionProvider),
			timeProvider,
		)
		productService = productsvc.New(
			pgproduct.New(transactionProvider),
		)
		cartService = cartsvc.New(
			cartRedis,
		)
		orderService = ordersvc.New(
			transactionProvider,
			pgorder.New(driver, transactionProvider),
			timeProvider,
		)
		currencyService = currencysvc.New(
			pgcurrency.New(transactionProvider),
		)
		notificationService = notificationsvc.New(
			pgoutbox.New(transactionProvider),
			sender,
		)
		cartOrderUsecase = cartorder.New(
			transactionProvider,
			storeService,
			productService,
			cartService,
			orderService,
			currencyService,
			notificationService,
		)
		historyOrderUsecase = history.New(
			cfg.OrderHistory.PageSize,
			orderService,
			currencyService,
			notificationService,
		)
		paymentOrderUsecase = payment.New(
			transactionProvider,
			storeService,
			orderService,
			productService,
			currencyService,
			dailyPosition,
			verificationcodegenerator.New(),
			qrcodegenerator.New(),
			timeProvider,
			notificationService,
		)
	)

	if err := app.Start(
		ctx,
		cfg.Bot,
		cartOrderUsecase,
		nil,
		historyOrderUsecase,
		paymentOrderUsecase,
		buttonRepository,
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
