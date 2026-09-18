package messagesenderobsolete

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) SendPNGMarkdown(
	ctx context.Context,
	chatID msginfo.ChatID,
	caption string,
	png []byte,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: chatID.Int64(),
		Photo: &models.InputFileUpload{
			Data: bytes.NewReader(png),
		},
		Caption:     caption,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send photo: %w", err)
	}

	return nil
}
