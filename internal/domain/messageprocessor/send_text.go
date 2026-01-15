package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) SendTextPlain(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, Message{
		ChatID:  chatID,
		Text:    text,
		Type:    MessageTypePlain,
		Buttons: rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (m *MessageProcessor) SendTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, Message{
		ChatID:  chatID,
		Text:    text,
		Type:    MessageTypeMarkdown,
		Buttons: rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
