package messagesenderobsolete

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	if _, err := m.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID.Int64(),
		MessageID: messageID.Int(),
	}); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
