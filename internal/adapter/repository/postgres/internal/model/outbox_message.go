package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type MessageType string

const (
	MessageTypePlain    MessageType = "plain"
	MessageTypeMarkdown MessageType = "markdown"
	MessageTypePNG      MessageType = "png"
)

type OutboxMessage struct {
	ID             int           `db:"id"`
	ChatID         int64         `db:"chat_id"`
	ReplyMessageID sql.NullInt64 `db:"reply_msg_id"`
	Text           string        `db:"msg_text"`
	Type           MessageType   `db:"msg_type"`
	Payload        []byte        `db:"payload"`
	Button         jsonb.JSONB   `db:"buttons"`
	IsDispatched   bool          `db:"is_dispatched"`
	CreatedAt      time.Time     `db:"created_at"`
	DispatchedAt   sql.NullTime  `db:"dispatched_at"`
}

func ToDBMessageType(msgType messageprocessor.MessageTextType) MessageType {
	switch msgType {
	case messageprocessor.MessageTextTypePlain:
		return MessageTypePlain
	case messageprocessor.MessageTextTypeMarkdown:
		return MessageTypeMarkdown
	}

	return ""
}

func ToMessageType(mt MessageType) messageprocessor.MessageTextType {
	switch mt {
	case MessageTypePlain:
		return messageprocessor.MessageTextTypePlain

	case MessageTypeMarkdown:
		return messageprocessor.MessageTextTypeMarkdown

	case MessageTypePNG:
		return 0
	}

	return 0
}

func ToOutboxMessage(msg OutboxMessage) outboxprocessor.OutboxMessage {
	return outboxprocessor.OutboxMessage{
		ID: msg.ID,
		Message: messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(msg.ChatID),
			ReplyMsgID: msginfo.MessageIDFromInt(int(msg.ReplyMessageID.Int64)),
			Text:       msg.Text,
			Type:       ToMessageType(msg.Type),
			Payload:    msg.Payload,
		},
	}
}

func ToOutboxMessages(dbMsgs []OutboxMessage) []outboxprocessor.OutboxMessage {
	if len(dbMsgs) == 0 {
		return nil
	}

	outboxMsgs := make([]outboxprocessor.OutboxMessage, 0, len(dbMsgs))

	for _, m := range dbMsgs {
		outboxMsgs = append(outboxMsgs, ToOutboxMessage(m))
	}

	return outboxMsgs
}
