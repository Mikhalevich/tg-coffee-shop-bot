package messagesvc

import (
	"context"
	"fmt"
)

func (s *Service) AnswerPayment(
	ctx context.Context,
	paymentID string,
	success bool,
	errorMsg string,
) error {
	if err := s.sender.AnswerPayment(
		ctx,
		paymentID,
		success,
		errorMsg,
	); err != nil {
		return fmt.Errorf("answer payment: %w", err)
	}

	return nil
}
