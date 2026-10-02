package notificationsvc

import (
	"context"
	"fmt"
)

func (s *Service) PaymentAnswer(
	ctx context.Context,
	paymentID string,
	success bool,
) error {
	if err := s.sender.AnswerPayment(
		ctx,
		paymentID,
		success,
		paymentMsg(success),
	); err != nil {
		return fmt.Errorf("answer payment: %w", err)
	}

	return nil
}

func paymentMsg(success bool) string {
	if success {
		return "Payment success"
	}

	return "Payment failure"
}
