package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (t *TGHandler) GetActiveOrder(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if err := t.viewActiveOrderUsecase.ViewActiveOrder(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("get active order: %w", err)
	}

	return nil
}
