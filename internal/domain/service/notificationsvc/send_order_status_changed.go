package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (s *Service) SendOrderStatusChanged(
	ctx context.Context,
	chatID msginfo.ChatID,
	newStatus order.Status,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   s.makeChangedOrderStatusMarkdownMsg(newStatus),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (s *Service) makeChangedOrderStatusMarkdownMsg(newStatus order.Status) string {
	return fmt.Sprintf("your order status changed to *%s*",
		s.escaper.EscapeMarkdown(newStatus.HumanReadable()),
	)
}
