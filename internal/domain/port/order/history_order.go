package order

import (
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
)

type HistoryOrder struct {
	ID           ID
	SerialNumber int
	Status       Status
	CurrencyID   currency.ID
	CreatedAt    time.Time
	TotalPrice   int
}
type Page struct {
	Number int
	Total  int
}

func (p Page) HasNext() bool {
	return p.Number < p.Total
}

func (p Page) Next() int {
	return p.Number + 1
}

func (p Page) HasPrevious() bool {
	return p.Number > 1
}

func (p Page) Previous() int {
	return p.Number - 1
}
