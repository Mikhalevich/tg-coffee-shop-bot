package notificationsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type Service struct {
	sender  Sender
	escaper MarkdownEscaper
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
) *Service {
	return &Service{
		sender:  sender,
		escaper: escaper,
	}
}
