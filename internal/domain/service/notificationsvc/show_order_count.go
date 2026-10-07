package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) ShowOrderCount(
	ctx context.Context,
	chatID msginfo.ChatID,
	count int,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   fmt.Sprintf("*%d*", count),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
