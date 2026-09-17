package messagesvc

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type SenderMessage struct {
	ChatID     msginfo.ChatID
	ReplyMsgID msginfo.MessageID
	Text       string
	Type       msginfo.MessageType
	Payload    []byte
	Buttons    []button.InlineKeyboardButtonRow
}

type SenderInvoice struct {
	ChatID      msginfo.ChatID
	Title       string
	Description string
	Currency    string
	Payload     string
	Labels      []msginfo.LabeledPrice
	Buttons     []button.InlineKeyboardButtonRow
}

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg SenderMessage,
	) error
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
	SendInvoice(
		ctx context.Context,
		invoice SenderInvoice,
	) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type ButtonRepository interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) error
}

type Service struct {
	sender           Sender
	escaper          MarkdownEscaper
	buttonRepository ButtonRepository
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
) *Service {
	return &Service{
		sender:           sender,
		escaper:          escaper,
		buttonRepository: buttonRepository,
	}
}
