package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) CartOrderUnavailable(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   message.CartOrderUnavailable(),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
