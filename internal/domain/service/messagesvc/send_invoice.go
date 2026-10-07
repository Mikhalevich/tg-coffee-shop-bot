package messagesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) SendInvoice(
	ctx context.Context,
	invoice order.Invoice,
) error {
	inlineButtons, err := s.SetButtonRows(ctx, invoice.Buttons...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := s.sender.SendInvoice(ctx, SenderInvoice{
		ChatID:       invoice.ChatID,
		Title:        invoice.Title,
		Description:  invoice.Description,
		CurrencyCode: invoice.CurrencyCode,
		OrderID:      invoice.OrderID,
		Labels:       invoice.Labels,
		Buttons:      inlineButtons,
	}); err != nil {
		return fmt.Errorf("sender send message: %w", err)
	}

	return nil
}
