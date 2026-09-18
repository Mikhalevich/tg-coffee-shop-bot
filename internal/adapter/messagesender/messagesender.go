package messagesender

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/messagesvc"
)

var (
	_ messagesvc.Sender          = (*messageSender)(nil)
	_ messagesvc.MarkdownEscaper = (*messageSender)(nil)
)

type messageSender struct {
	bot          *bot.Bot
	paymentToken string
}

func New(bot *bot.Bot, paymentToken string) *messageSender {
	return &messageSender{
		bot:          bot,
		paymentToken: paymentToken,
	}
}

func makeButtonsMarkup(rows ...button.InlineKeyboardButtonRow) models.ReplyMarkup {
	if len(rows) == 0 {
		return nil
	}

	keyboard := make([][]models.InlineKeyboardButton, 0, len(rows))

	for _, row := range rows {
		buttonRow := make([]models.InlineKeyboardButton, 0, len(row))

		for _, b := range row {
			buttonRow = append(buttonRow, models.InlineKeyboardButton{
				Text:         b.Caption,
				CallbackData: b.ID.String(),
				Style:        b.Style.String(),
				URL:          b.URL,
			})
		}

		keyboard = append(keyboard, buttonRow)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
