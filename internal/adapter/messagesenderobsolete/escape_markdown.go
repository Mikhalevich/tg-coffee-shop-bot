package messagesenderobsolete

import (
	"github.com/go-telegram/bot"
)

func (m *messageSender) EscapeMarkdown(s string) string {
	return bot.EscapeMarkdown(s)
}
