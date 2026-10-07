package event

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type OutboxMessage struct {
	ID             int                 `json:"id"`
	ChatID         int64               `json:"chat_id"`
	ReplyMessageID *int64              `json:"reply_msg_id"`
	MessageText    string              `json:"msg_text"`
	MessageType    msginfo.MessageType `json:"msg_type"`
	Payload        *string             `json:"payload"`
	Buttons        string              `json:"buttons"`
}
