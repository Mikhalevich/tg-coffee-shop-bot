package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
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

func intToNullInt64(value int) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(value),
		Valid: true,
	}
}

func ToDBOutboxMessage(msg messageprocessor.Message) (OutboxMessage, error) {
	jbButtons, err := jsonbFromButtonRows(msg.Buttons)
	if err != nil {
		return OutboxMessage{}, fmt.Errorf("jsonb from buttons: %w", err)
	}

	return OutboxMessage{
		ChatID:         msg.ChatID.Int64(),
		ReplyMessageID: intToNullInt64(msg.ReplyMsgID.Int()),
		Text:           msg.Text,
		Type:           ToDBMessageType(msg.Type),
		Payload:        msg.Payload,
		Button:         jbButtons,
	}, nil
}

func jsonbFromButtonRows(buttons []button.ButtonRow) (jsonb.JSONB, error) {
	if buttons == nil {
		return jsonb.NewString("[]"), nil
	}

	jbButtons, err := jsonb.NewFromMarshaler(buttons)
	if err != nil {
		return jsonb.NewNull(), fmt.Errorf("jsonb marshaler: %w", err)
	}

	return jbButtons, nil
}

func ToDBMessageType(msgType messageprocessor.MessageType) MessageType {
	switch msgType {
	case messageprocessor.MessageTypePlain:
		return MessageTypePlain
	case messageprocessor.MessageTypeMarkdown:
		return MessageTypeMarkdown
	case messageprocessor.MessageTypePNG:
		return MessageTypePNG
	}

	return ""
}

func ToMessageType(mt MessageType) messageprocessor.MessageType {
	switch mt {
	case MessageTypePlain:
		return messageprocessor.MessageTypePlain

	case MessageTypeMarkdown:
		return messageprocessor.MessageTypeMarkdown

	case MessageTypePNG:
		return messageprocessor.MessageTypePNG
	}

	return 0
}

func ToOutboxMessage(msg OutboxMessage) (outboxprocessor.OutboxMessage, error) {
	var buttons []button.ButtonRow
	if err := jsonb.ConvertTo(msg.Button, &buttons); err != nil {
		return outboxprocessor.OutboxMessage{}, fmt.Errorf("convert jsonb to button rows: %w", err)
	}

	return outboxprocessor.OutboxMessage{
		ID: msg.ID,
		Message: messageprocessor.Message{
			ChatID:     msginfo.ChatIDFromInt(msg.ChatID),
			ReplyMsgID: msginfo.MessageIDFromInt(int(msg.ReplyMessageID.Int64)),
			Text:       msg.Text,
			Type:       ToMessageType(msg.Type),
			Payload:    msg.Payload,
		},
	}, nil
}

func ToOutboxMessages(dbMsgs []OutboxMessage) ([]outboxprocessor.OutboxMessage, error) {
	if len(dbMsgs) == 0 {
		return nil, nil
	}

	outboxMsgs := make([]outboxprocessor.OutboxMessage, 0, len(dbMsgs))

	for _, m := range dbMsgs {
		outboxMsg, err := ToOutboxMessage(m)
		if err != nil {
			return nil, fmt.Errorf("make outbox message: %w", err)
		}

		outboxMsgs = append(outboxMsgs, outboxMsg)
	}

	return outboxMsgs, nil
}
