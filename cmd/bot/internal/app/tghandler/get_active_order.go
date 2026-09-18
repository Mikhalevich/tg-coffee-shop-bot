package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) GetActiveOrder(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.actionProcessor.GetActiveOrder(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt64(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("get active order: %w", err)
	}

	return nil
}
