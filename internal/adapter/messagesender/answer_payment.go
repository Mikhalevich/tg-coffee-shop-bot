package messagesender

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func (m *messageSender) AnswerPayment(
	ctx context.Context,
	paymentID string,
	success bool,
	errorMsg string,
) error {
	if _, err := m.bot.AnswerPreCheckoutQuery(
		ctx,
		&bot.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: paymentID,
			OK:                 success,
			ErrorMessage:       errorMsg,
		},
	); err != nil {
		return fmt.Errorf("answer precheckout query: %w", err)
	}

	return nil
}
