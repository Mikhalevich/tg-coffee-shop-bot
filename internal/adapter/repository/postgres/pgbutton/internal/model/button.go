package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type Button struct {
	ID                   string `db:"id"`
	Caption              string `db:"caption"`
	Operation            string `db:"operation"`
	IsDeleteAfterProcess bool   `db:"is_delete_after_process"`
	Style                string `db:"style"`
	URL                  string `db:"url"`
	Payload              []byte `db:"payload"`
	Pay                  bool   `db:"pay"`
}

func (b Button) ToDom() button.Button {
	return button.Button{
		ID:                   button.IDFromString(b.ID),
		Caption:              b.Caption,
		Operation:            button.Operation(b.Operation),
		IsDeleteAfterProcess: b.IsDeleteAfterProcess,
		Style:                button.StyleFromString(b.Style),
		URL:                  b.URL,
		Payload:              b.Payload,
		Pay:                  b.Pay,
	}
}

// ToDBButtons flattens button rows skipping pay buttons:
// they have no id and are never returned as callbacks.
func ToDBButtons(rows ...button.ButtonRow) []Button {
	var dbButtons []Button

	for _, row := range rows {
		for _, btn := range row {
			if btn.Pay {
				continue
			}

			dbButtons = append(dbButtons, Button{
				ID:                   btn.ID.String(),
				Caption:              btn.Caption,
				Operation:            btn.Operation.String(),
				IsDeleteAfterProcess: btn.IsDeleteAfterProcess,
				Style:                btn.Style.String(),
				URL:                  btn.URL,
				Payload:              btn.Payload,
				Pay:                  btn.Pay,
			})
		}
	}

	return dbButtons
}
