package button

type InlineKeyboardButton struct {
	ID      ID
	Caption string
	Style   Style
	URL     string
}

func ToInlineKeyboardButton(b Button) InlineKeyboardButton {
	return InlineKeyboardButton{
		ID:      b.ID,
		Caption: b.Caption,
		Style:   b.Style,
		URL:     b.URL,
	}
}

type InlineKeyboardButtonRow []InlineKeyboardButton

func InlineRow(buttons ...InlineKeyboardButton) InlineKeyboardButtonRow {
	return buttons
}
