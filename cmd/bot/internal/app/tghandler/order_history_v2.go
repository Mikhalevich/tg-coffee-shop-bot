package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) OrderHistoryV2(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.historyProcessorV2.Show(
		ctx,
		msginfo.Info{
			ChatID:    msginfo.ChatIDFromInt(msg.ChatID),
			MessageID: msginfo.MessageIDFromInt(msg.MessageID),
		},
	); err != nil {
		return fmt.Errorf("history orders first: %w", err)
	}

	return nil
}
