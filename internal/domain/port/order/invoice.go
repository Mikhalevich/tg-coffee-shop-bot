package order

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
)

type LabeledPrice struct {
	Label  string
	Amount int
}

type Invoice struct {
	ChatID       msginfo.ChatID
	Title        string
	Description  string
	CurrencyCode string
	OrderID      ID
	Labels       []LabeledPrice
	Buttons      []button.ButtonRow
}
