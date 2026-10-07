package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/null"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
)

type MessageType string

const (
	MessageTypePlain    MessageType = "plain"
	MessageTypeMarkdown MessageType = "markdown"
	MessageTypePNG      MessageType = "png"
)

type Message struct {
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

func ToDBOutboxMessage(msg msginfo.Message) (Message, error) {
	jbButtons, err := jsonbFromSlice(msg.Buttons)
	if err != nil {
		return Message{}, fmt.Errorf("jsonb from buttons: %w", err)
	}

	return Message{
		ChatID:         msg.ChatID.Int64(),
		ReplyMessageID: null.Int64Positive(int64(msg.ReplyMsgID.Int())),
		Text:           msg.Text,
		Type:           ToDBMessageType(msg.Type),
		Payload:        msg.Payload,
		Button:         jbButtons,
	}, nil
}

func jsonbFromSlice[T any](elements []T) (jsonb.JSONB, error) {
	if len(elements) == 0 {
		return jsonb.NewString("[]"), nil
	}

	jbButtons, err := jsonb.NewFromMarshaler(elements)
	if err != nil {
		return jsonb.NewNull(), fmt.Errorf("jsonb marshaler: %w", err)
	}

	return jbButtons, nil
}

func ToDBMessageType(msgType msginfo.MessageType) MessageType {
	switch msgType {
	case msginfo.MessageTypePlain:
		return MessageTypePlain
	case msginfo.MessageTypeMarkdown:
		return MessageTypeMarkdown
	case msginfo.MessageTypePNG:
		return MessageTypePNG
	}

	return ""
}

func ToMessageType(mt MessageType) msginfo.MessageType {
	switch mt {
	case MessageTypePlain:
		return msginfo.MessageTypePlain

	case MessageTypeMarkdown:
		return msginfo.MessageTypeMarkdown

	case MessageTypePNG:
		return msginfo.MessageTypePNG
	}

	return 0
}

func ToOutboxMessage(msg Message) (outboxmsg.Message, error) {
	var buttons []button.ButtonRow
	if err := jsonb.ConvertTo(msg.Button, &buttons); err != nil {
		return outboxmsg.Message{}, fmt.Errorf("convert jsonb to button rows: %w", err)
	}

	return outboxmsg.Message{
		ID: msg.ID,
		Message: msginfo.Message{
			ChatID:     msginfo.ChatIDFromInt64(msg.ChatID),
			ReplyMsgID: msginfo.MessageIDFromInt(int(msg.ReplyMessageID.Int64)),
			Text:       msg.Text,
			Type:       ToMessageType(msg.Type),
			Payload:    msg.Payload,
		},
	}, nil
}

func ToOutboxMessages(dbMsgs []Message) ([]outboxmsg.Message, error) {
	if len(dbMsgs) == 0 {
		return nil, nil
	}

	outboxMsgs := make([]outboxmsg.Message, 0, len(dbMsgs))

	for _, m := range dbMsgs {
		outboxMsg, err := ToOutboxMessage(m)
		if err != nil {
			return nil, fmt.Errorf("make outbox message: %w", err)
		}

		outboxMsgs = append(outboxMsgs, outboxMsg)
	}

	return outboxMsgs, nil
}
