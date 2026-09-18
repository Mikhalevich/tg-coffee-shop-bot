package notificationsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/orderbyid"
)

var (
	_ cartorder.NotificationService = (*Service)(nil)
	_ orderbyid.NotificationService = (*Service)(nil)
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
	SendInvoice(ctx context.Context, invoice msginfo.Invoice) error
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
