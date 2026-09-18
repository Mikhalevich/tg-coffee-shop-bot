package messagesenderobsolete

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func (m *messageSender) AnswerOrderPayment(ctx context.Context, paymentID string, ok bool, errorMsg string) error {
	if _, err := m.bot.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: paymentID,
		OK:                 ok,
		ErrorMessage:       errorMsg,
	}); err != nil {
		return fmt.Errorf("answer precheckout query: %w", err)
	}

	return nil
}
