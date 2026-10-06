package model

import (
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/outboxmsg"
)

type Invoice struct {
	ID           int         `db:"id"`
	ChatID       int64       `db:"chat_id"`
	Title        string      `db:"title"`
	Description  string      `db:"description"`
	CurrencyCode string      `db:"currency_code"`
	OrderID      int         `db:"order_id"`
	Labels       jsonb.JSONB `db:"labels"`
	Buttons      jsonb.JSONB `db:"buttons"`
}

func ToDBInvoice(invoice order.Invoice) (Invoice, error) {
	jbLabels, err := jsonbFromSlice(invoice.Labels)
	if err != nil {
		return Invoice{}, fmt.Errorf("jsonb from labels: %w", err)
	}

	jbButtons, err := jsonbFromSlice(invoice.Buttons)
	if err != nil {
		return Invoice{}, fmt.Errorf("jsonb from buttons: %w", err)
	}

	return Invoice{
		ChatID:       invoice.ChatID.Int64(),
		Title:        invoice.Title,
		Description:  invoice.Description,
		CurrencyCode: invoice.CurrencyCode,
		OrderID:      invoice.OrderID.Int(),
		Labels:       jbLabels,
		Buttons:      jbButtons,
	}, nil
}

func (i Invoice) ToOutbox() (outboxmsg.Invoice, error) {
	var labels []order.LabeledPrice
	if err := jsonb.ConvertTo(i.Labels, &labels); err != nil {
		return outboxmsg.Invoice{}, fmt.Errorf("convert jsonb to labels: %w", err)
	}

	var buttons []button.ButtonRow
	if err := jsonb.ConvertTo(i.Buttons, &buttons); err != nil {
		return outboxmsg.Invoice{}, fmt.Errorf("convert jsonb to button rows: %w", err)
	}

	return outboxmsg.Invoice{
		ID: i.ID,
		Invoice: order.Invoice{
			ChatID:       msginfo.ChatIDFromInt64(i.ChatID),
			Title:        i.Title,
			Description:  i.Description,
			CurrencyCode: i.CurrencyCode,
			OrderID:      order.IDFromInt(i.OrderID),
			Labels:       labels,
			Buttons:      buttons,
		},
	}, nil
}

func ToOutboxInvoices(dbInvoices []Invoice) ([]outboxmsg.Invoice, error) {
	domInvoices := make([]outboxmsg.Invoice, 0, len(dbInvoices))

	for _, di := range dbInvoices {
		domInvoice, err := di.ToOutbox()
		if err != nil {
			return nil, fmt.Errorf("convert to dom invoice: %w", err)
		}

		domInvoices = append(domInvoices, domInvoice)
	}

	return domInvoices, nil
}
