package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

func (s *Service) PaymentConfirmed(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord order.Order,
	curr currency.Currency,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
	qrCodeImage []byte,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePNG,
			Text: s.formatOrder(
				ord,
				productsInfo,
				curr,
				queuePosition,
			),
			Payload: qrCodeImage,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
