package messagesenderobsolete

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) ReplyText(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("reply message plain: %w", err)
	}

	return nil
}

func (m *messageSender) ReplyTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyToMsgID msginfo.MessageID,
	text string,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID.Int64(),
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID.Int(),
		},
		ParseMode:   models.ParseModeMarkdown,
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("reply message markdown: %w", err)
	}

	return nil
}
