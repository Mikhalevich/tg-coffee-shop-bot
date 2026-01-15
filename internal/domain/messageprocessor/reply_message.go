package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *MessageProcessor) ReplyTextPlain(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, Message{
		ChatID:     chatID,
		ReplyMsgID: replyMessageID,
		Text:       text,
		Type:       MessageTypePlain,
		Buttons:    rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (m *MessageProcessor) ReplyTextMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	replyMessageID msginfo.MessageID,
	text string,
	rows ...button.ButtonRow,
) error {
	if err := m.SendMessage(ctx, Message{
		ChatID:     chatID,
		ReplyMsgID: replyMessageID,
		Text:       text,
		Type:       MessageTypeMarkdown,
		Buttons:    rows,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
