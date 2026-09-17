package msginfo

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type MessageID int

func (m MessageID) Int() int {
	return int(m)
}

func MessageIDFromInt(id int) MessageID {
	return MessageID(id)
}

type ChatID int64

func (c ChatID) Int64() int64 {
	return int64(c)
}

func ChatIDFromInt(id int64) ChatID {
	return ChatID(id)
}

type Info struct {
	ChatID    ChatID
	MessageID MessageID
}

type MessageType int

const (
	MessageTypePlain MessageType = iota + 1
	MessageTypeMarkdown
	MessageTypePNG
)

func (mt MessageType) Int() int {
	return int(mt)
}

func MessageTypeFromInt(t int) MessageType {
	return MessageType(t)
}

type Message struct {
	ChatID       ChatID
	ReplyMsgID   MessageID
	Text         string
	Type         MessageType
	Payload      []byte
	Buttons      []button.ButtonRow
	VisibilityAt time.Time
}

type LabeledPrice struct {
	Label  string
	Amount int
}

type Invoice struct {
	ChatID      ChatID
	Title       string
	Description string
	Currency    string
	Payload     string
	Labels      []LabeledPrice
	Buttons     []button.ButtonRow
}
