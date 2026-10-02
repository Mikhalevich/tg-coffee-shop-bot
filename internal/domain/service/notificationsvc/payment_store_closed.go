package notificationsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
)

func (s *Service) PaymentStoreClosed(
	ctx context.Context,
	paymentID string,
	currentTime time.Time,
	nextWorkingTime time.Time,
) error {
	if err := s.sender.AnswerPayment(
		ctx,
		paymentID,
		false,
		message.StoreClosed(currentTime, nextWorkingTime),
	); err != nil {
		return fmt.Errorf("answer payment: %w", err)
	}

	return nil
}
