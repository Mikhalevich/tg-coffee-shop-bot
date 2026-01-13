package messageprocessor

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
)

type MessageTextType int

const (
	MessageTextTypePlain MessageTextType = iota + 1
	MessageTextTypeMarkdown
)

type Sender interface {
	SendText(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	SendTextMarkdown(
		ctx context.Context,
		chatID msginfo.ChatID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	ReplyText(
		ctx context.Context,
		chatID msginfo.ChatID,
		replyToMsgID msginfo.MessageID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	ReplyTextMarkdown(
		ctx context.Context,
		chatID msginfo.ChatID,
		replyToMsgID msginfo.MessageID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	SendPNGMarkdown(
		ctx context.Context,
		chatID msginfo.ChatID,
		caption string,
		png []byte,
		rows ...button.InlineKeyboardButtonRow,
	) error
	EditText(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
	SendOrderInvoice(
		ctx context.Context,
		chatID msginfo.ChatID,
		title string,
		description string,
		ord *order.Order,
		productsInfo map[product.ProductID]product.Product,
		currency string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	AnswerOrderPayment(
		ctx context.Context,
		paymentID string,
		ok bool,
		errorMsg string,
	) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type ButtonRepository interface {
	SetButton(ctx context.Context, btn button.Button) error
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) error

	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	IsNotFoundError(err error) bool
}

type Message struct {
	ChatID     msginfo.ChatID
	ReplyMsgID msginfo.MessageID
	Text       string
	Type       MessageTextType
	Payload    []byte
	Buttons    []button.ButtonRow
}

type MessageProcessor struct {
	sender           Sender
	escaper          MarkdownEscaper
	buttonRepository ButtonRepository
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
) *MessageProcessor {
	return &MessageProcessor{
		sender:           sender,
		escaper:          escaper,
		buttonRepository: buttonRepository,
	}
}
