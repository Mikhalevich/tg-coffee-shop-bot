package messagesender

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/messagesvc"
)

//nolint:funlen
func (m *messageSender) SendMessage(
	ctx context.Context,
	msg messagesvc.SenderMessage,
) error {
	switch msg.Type {
	case msginfo.MessageTypePlain, msginfo.MessageTypeMarkdown:
		if msg.ReplyMsgID.IsValid() {
			if _, err := m.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      msg.ChatID.Int64(),
				MessageID:   msg.ReplyMsgID.Int(),
				Text:        msg.Text,
				ParseMode:   textParseMode(msg.Type),
				ReplyMarkup: makeButtonsMarkup(msg.Buttons...),
			}); err != nil {
				return fmt.Errorf("edit message text: %w", err)
			}

			break
		}

		if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          msg.ChatID.Int64(),
			Text:            msg.Text,
			ParseMode:       textParseMode(msg.Type),
			ReplyParameters: replyParameters(msg.ReplyMsgID),
			ReplyMarkup:     makeButtonsMarkup(msg.Buttons...),
		}); err != nil {
			return fmt.Errorf("send text message: %w", err)
		}

	case msginfo.MessageTypePNG:
		if msg.ReplyMsgID.IsValid() {
			if _, err := m.bot.EditMessageMedia(ctx, &bot.EditMessageMediaParams{
				ChatID:    msg.ChatID.Int64(),
				MessageID: msg.ReplyMsgID.Int(),
				Media: &models.InputMediaPhoto{
					Media:           "attach://filename",
					Caption:         msg.Text,
					ParseMode:       models.ParseModeMarkdown,
					MediaAttachment: bytes.NewReader(msg.Payload),
				},
				ReplyMarkup: makeButtonsMarkup(msg.Buttons...),
			}); err != nil {
				return fmt.Errorf("edit media: %w", err)
			}

			break
		}

		if _, err := m.bot.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: msg.ChatID.Int64(),
			Photo: &models.InputFileUpload{
				Data: bytes.NewReader(msg.Payload),
			},
			Caption:     msg.Text,
			ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: makeButtonsMarkup(msg.Buttons...),
		}); err != nil {
			return fmt.Errorf("send photo: %w", err)
		}

	default:
		return fmt.Errorf("invalid message type: %v", msg.Type)
	}

	return nil
}

func textParseMode(mt msginfo.MessageType) models.ParseMode {
	if mt == msginfo.MessageTypeMarkdown {
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
