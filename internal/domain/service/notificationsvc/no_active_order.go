package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) NoActiveOrder(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   "No active orders",
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
