package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) SendText(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send text plain: %w", err)
	}

	return nil
}

func (m *messageSender) SendTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID.Int64(),
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send text markdown: %w", err)
	}

	return nil
}
