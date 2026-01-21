package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/outboxprocessor"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

type OutboxInvoice struct {
	ID      int    `db:"id"`
	ChatID  int64  `db:"chat_id"`
	Text    string `db:"msg_text"`
	OrderID int    `db:"order_id"`
}

func ToOutboxInvoice(dbInvoice OutboxInvoice) outboxprocessor.OutboxInvoice {
	return outboxprocessor.OutboxInvoice{
		ID:      dbInvoice.ID,
		ChatID:  msginfo.ChatIDFromInt(dbInvoice.ChatID),
		Text:    dbInvoice.Text,
		OrderID: order.IDFromInt(dbInvoice.OrderID),
	}
}

func ToOutboxInvoices(dbInvoices []OutboxInvoice) []outboxprocessor.OutboxInvoice {
	invoices := make([]outboxprocessor.OutboxInvoice, 0, len(dbInvoices))

	for _, di := range dbInvoices {
		invoices = append(invoices, ToOutboxInvoice(di))
	}

	return invoices
}
