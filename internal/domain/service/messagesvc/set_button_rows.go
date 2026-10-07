package messagesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

func (s *Service) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) ([]button.InlineKeyboardButtonRow, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	if err := s.buttonRepository.SetButtonRows(ctx, rows...); err != nil {
		return nil, fmt.Errorf("set button rows: %w", err)
	}

	inlineButtonRows := make([]button.InlineKeyboardButtonRow, 0, len(rows))

	for _, row := range rows {
		buttonRow := make([]button.InlineKeyboardButton, 0, len(row))

		for _, btn := range row {
			buttonRow = append(buttonRow, button.ToInlineKeyboardButton(btn))
		}

		if len(buttonRow) > 0 {
			inlineButtonRows = append(inlineButtonRows, buttonRow)
		}
	}

	return inlineButtonRows, nil
}
