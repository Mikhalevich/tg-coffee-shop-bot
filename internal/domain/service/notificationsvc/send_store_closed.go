package notificationsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

func (s *Service) SendStoreClosed(
	ctx context.Context,
	chatID msginfo.ChatID,
	currentTime time.Time,
	nextWorkingTime time.Time,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   storeClosedMsg(currentTime, nextWorkingTime),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func storeClosedMsg(currentTime, openTime time.Time) string {
	return fmt.Sprintf("Closed. Will be opened after %s at %s",
		openTime.Sub(currentTime).Truncate(time.Minute).String(),
		openTime.Format("Monday 15:04 MST"))
}
