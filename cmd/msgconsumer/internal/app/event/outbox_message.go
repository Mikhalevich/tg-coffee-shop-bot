package event

import "github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"

type OutboxMessageType string

const (
	OutboxMessageTypePlain    OutboxMessageType = "plain"
	OutboxMessageTypeMarkdown OutboxMessageType = "markdown"
	OutboxMessageTypePNG      OutboxMessageType = "png"
)

type OutboxMessage struct {
	ID             int               `json:"id"`
	ChatID         int64             `json:"chat_id"`
	ReplyMessageID *int64            `json:"reply_msg_id"`
	MessageText    string            `json:"msg_text"`
	MessageType    OutboxMessageType `json:"msg_type"`
	Payload        *string           `json:"payload"`
	Buttons        string            `json:"buttons"`
}

func ToMessageType(mt OutboxMessageType) messageprocessor.MessageType {
	switch mt {
	case OutboxMessageTypePlain:
		return messageprocessor.MessageTypePlain
	case OutboxMessageTypeMarkdown:
		return messageprocessor.MessageTypeMarkdown
	case OutboxMessageTypePNG:
		return messageprocessor.MessageTypePNG
	}

	return 0
}
