package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app/kafkaconsumer"
	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/buttonrespository"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
)

func StartConsumer(
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

	var (
		consumer     = kafkaconsumer.New(cfg.Kafka)
		sender       = messagesender.New(botAPI, cfg.Bot.PaymentToken)
		msgProcessor = messageprocessor.New(sender, sender, buttonRepository)
	)

	if err := app.New(consumer, msgProcessor).Start(ctx); err != nil {
		return fmt.Errorf("start app: %w", err)
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
