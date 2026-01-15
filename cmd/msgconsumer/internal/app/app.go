package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/app/event"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type Consumer interface {
	Consume(
		ctx context.Context,
		processFn func(ctx context.Context, payload []byte) error,
	) error
}

type MessageSender interface {
	SendMessage(
		ctx context.Context,
		msg messageprocessor.Message,
	) error
}

type App struct {
	consumer Consumer
	sender   MessageSender
}

func New(
	consumer Consumer,
	sender MessageSender,
) *App {
	return &App{
		consumer: consumer,
		sender:   sender,
	}
}

func (a *App) Start(ctx context.Context) error {
	if err := a.consumer.Consume(ctx, func(ctx context.Context, payload []byte) error {
		var msg event.OutboxMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return fmt.Errorf("unmarshal message: %w", err)
		}

		if err := a.sender.SendMessage(ctx, messageprocessor.Message{
			ChatID: msginfo.ChatIDFromInt(msg.ChatID),
			Text:   msg.MessageText,
			Type:   event.ToMessageType(msg.MessageType),
		}); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return fmt.Errorf("unknown message type: %s", msg.MessageType)
	}); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}
