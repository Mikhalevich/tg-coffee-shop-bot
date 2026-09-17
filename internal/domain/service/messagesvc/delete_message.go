package messagesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	if messageID.Int() == 0 {
		return nil
	}

	if err := s.sender.DeleteMessage(ctx, chatID, messageID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
