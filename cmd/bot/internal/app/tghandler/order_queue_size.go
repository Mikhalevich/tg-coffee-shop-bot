package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) OrderQueueSize(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.actionProcessor.QueueSize(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt64(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("order queue size: %w", err)
	}

	return nil
}
