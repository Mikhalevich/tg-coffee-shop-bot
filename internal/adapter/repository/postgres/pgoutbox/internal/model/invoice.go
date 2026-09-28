package model

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
)

type Invoice struct {
	ID      int    `db:"id"`
	ChatID  int64  `db:"chat_id"`
	Text    string `db:"msg_text"`
	OrderID int    `db:"order_id"`
}

func (i Invoice) ToDom() outboxmsg.Invoice {
	return outboxmsg.Invoice{
		ID:      i.ID,
		ChatID:  msginfo.ChatIDFromInt64(i.ChatID),
		Text:    i.Text,
		OrderID: order.IDFromInt(i.OrderID),
	}
}

func ToDomInvoices(dbInvoices []Invoice) []outboxmsg.Invoice {
	invoices := make([]outboxmsg.Invoice, 0, len(dbInvoices))

	for _, di := range dbInvoices {
		invoices = append(invoices, di.ToDom())
	}

	return invoices
}
