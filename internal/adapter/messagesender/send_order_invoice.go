package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (m *messageSender) SendOrderInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	title string,
	description string,
	currency string,
	payload string,
	labels []messageprocessor.LabeledPrice,
	rows ...button.InlineKeyboardButtonRow,
) error {
	if _, err := m.bot.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        chatID.Int64(),
		Title:         title,
		Description:   description,
		Payload:       payload,
		ProviderToken: m.paymentToken,
		Currency:      currency,
		Prices:        toLabeledPrices(labels),
		ReplyMarkup:   makeButtonsMarkup(rows...),
	}); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}

func toLabeledPrices(labels []messageprocessor.LabeledPrice) []models.LabeledPrice {
	prices := make([]models.LabeledPrice, 0, len(labels))

	for _, v := range labels {
		prices = append(prices, models.LabeledPrice{
			Label:  v.Label,
			Amount: v.Amount,
		})
	}

	return prices
}
