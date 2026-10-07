package notificationsvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/cartorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/activeorder"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/history"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/orderbyid"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/ordercancel"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/payment"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/customer/order/queuesize"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/nextpending"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/usecase/manager/order/updatestatus"
)

var (
	_ cartorder.NotificationService    = (*Service)(nil)
	_ orderbyid.NotificationService    = (*Service)(nil)
	_ activeorder.NotificationService  = (*Service)(nil)
	_ ordercancel.NotificationService  = (*Service)(nil)
	_ payment.NotificationService      = (*Service)(nil)
	_ history.NotificationService      = (*Service)(nil)
	_ nextpending.NotificationService  = (*Service)(nil)
	_ updatestatus.NotificationService = (*Service)(nil)
	_ queuesize.NotificationService    = (*Service)(nil)
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
	SendInvoice(ctx context.Context, invoice order.Invoice) error
	AnswerPayment(
		ctx context.Context,
		paymentID string,
		success bool,
		errorMsg string,
	) error
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
