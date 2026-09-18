package messagesenderobsolete

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) EditText(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID.Int64(),
		MessageID:   messageID.Int(),
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("eidt message text: %w", err)
	}

	return nil
}
