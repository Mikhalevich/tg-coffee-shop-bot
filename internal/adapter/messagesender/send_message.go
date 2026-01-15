package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) SendMessage(
	ctx context.Context,
	msg messageprocessor.SenderMessage,
) error {
	switch msg.Type {
	case messageprocessor.MessageTypePlain, messageprocessor.MessageTypeMarkdown:
		if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          msg.ChatID.Int64(),
			Text:            msg.Text,
			ParseMode:       parseMode(msg.Type),
			ReplyParameters: replyParameters(msg.ReplyMsgID),
			ReplyMarkup:     makeButtonsMarkup(msg.Buttons...),
		}); err != nil {
			return fmt.Errorf("send text message: %w", err)
		}

	case messageprocessor.MessageTypePNG:
		if err := m.SendPNGMarkdown(ctx, msg.ChatID, msg.Text, msg.Payload, msg.Buttons...); err != nil {
			return fmt.Errorf("send png: %w", err)
		}

	default:
		return fmt.Errorf("invalid message type: %v", msg.Type)
	}

	return nil
}

func parseMode(mt messageprocessor.MessageType) models.ParseMode {
	if mt == messageprocessor.MessageTypeMarkdown {
		return models.ParseModeMarkdown
	}

	return ""
}

func replyParameters(replyMsgID msginfo.MessageID) *models.ReplyParameters {
	if replyMsgID.Int() == 0 {
		return nil
	}

	return &models.ReplyParameters{
		MessageID: replyMsgID.Int(),
	}
}
