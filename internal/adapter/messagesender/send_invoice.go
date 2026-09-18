package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/service/messagesvc"
)

func (m *messageSender) SendInvoice(
	ctx context.Context,
	invoice messagesvc.SenderInvoice,
) error {
	if _, err := m.bot.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        invoice.ChatID.Int64(),
		Title:         invoice.Title,
		Description:   invoice.Description,
		Payload:       invoice.Payload,
		ProviderToken: m.paymentToken,
		Currency:      invoice.Currency,
		Prices:        toLabeledPrices(invoice.Labels),
		ReplyMarkup:   makeButtonsMarkup(invoice.Buttons...),
	}); err != nil {
		return fmt.Errorf("send invoice: %w", err)
	}

	return nil
}

func toLabeledPrices(labels []msginfo.LabeledPrice) []models.LabeledPrice {
	prices := make([]models.LabeledPrice, 0, len(labels))

	for _, v := range labels {
		prices = append(prices, models.LabeledPrice{
			Label:  v.Label,
			Amount: v.Amount,
		})
	}

	return prices
}
